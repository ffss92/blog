---
title: Leveraging Go generics for input validation
subtitle: |
  Go's generics may not be as extensive as those in other languages, but they still offer powerful capabilities. 
  In this article, I'll demonstrate how I utilize generics to simplify and enhance validation in my projects.
author: "@ffss"
draft: false
date: "2025-01-02"
tags:
  - Go
  - Generics
---

## Introduction - a simple pattern

In his excellent book, [Let's Go](https://lets-go.alexedwards.net/), [Alex Edwards](https://www.alexedwards.net)
introduces us a way to do validation in HTTP requests using a `Validator` struct that has a `map[string]string`
field for holding errors, which I'm going to refer as `problem` in this article.

### A basic validator package

Edwards starts by defining a `validator` package, which might look like this:

```go
package validator

type Validator struct {
    Problems map[string]string
}

// Adds a validation problem.
func (v *Validator) SetProblem(name, reason string) {
    if v.Problems == nil {
        v.Problems = make(map[string]string)
    }
    if _, ok := v.Problems[name]; !ok {
        v.Problems[name] = reason
    }
}

// Adds a validation problem if ok is false.
func (v *Validator) Check(ok bool, name, reason string) {
    if !ok {
        v.SetProblem(name, reason)
    }
}

// Checks if the validator has any problems.
func (v Validator) IsValid() bool {
    return len(v.Problems) == 0
}
```

He also provides helper functions to validate specific cases:

```go
package validator

// Checks if a string is not blank.
func NotBlank(value string) bool {
    return utf8.RuneCountInString(value) > 0
}

// Checks if a string has at least n chars.
func MinChars(value string, n int) bool {
    return utf8.RuneCountInString(value) >= n
}

// Checks if a string has at most n chars.
func MaxChars(value string, n int) bool {
    return utf8.RuneCountInString(value) <= n
}

// Checks if value is a valid email.
func IsEmail(value string) bool {
    addr, err := mail.ParseAddress(value)
    return err == nil && addr.Address == value
}

// Checks if a is equal to b.
func Equal[T comparable](a, b T) bool {
    return a == b
}


// And others...
```

### Applying validation

To apply validation, a `validator.Validator` struct is embedded in the request struct,
allowing us to easily validate the request, like this:

```go
package main

func handleCreateUser() http.HandlerFunc {
    type request struct {
        Name            string `json:"name"`
        Email           string `json:"email"`
        Password        string `json:"password"`

        validator.Validator `json:"-"` // Omit from decoding.
    }
    return func(w http.ResponseWriter, r *http.Request) {
        var req request
        err := json.NewDecoder()
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        req.Check(validator.NotBlank(req.Name), "name", "Required field")
        req.Check(validator.MaxChars(req.Name, 120), "name", "Must have a most 120 characters")
        req.Check(validator.NotBlank(req.Email), "email", "Required field")
        req.Check(validator.IsEmail(req.Email), "email", "Must be a valid email")
        req.Check(validator.MaxChars(req.Email, 120), "email", "Must have a most 120 characters")
        req.Check(validator.NotBlank(req.Password), "password", "Required field")
        req.Check(validator.MinChars(req.Password, 8), "password", "Must have at least 8 characters")
        req.Check(validator.MaxChars(req.Password, 72), "password", "Must have at most 72 characters")
        if !req.IsValid() {
            // Handle invalid request
            return
        }

        // Create new user
    }
}
```

### It's not a perfect solution

While functional, this method has its shortcomings:

1. **Verbosity**: Writing validation rules for each field can feel repetitive.
2. **Error-Prone**: Mistyped field names or inconsistent reasons can cause bugs.

A naive approach would be updating our helper validation functions to return `bool, string` instead of just `bool`,
but so far we only validated `string`s, what if we wanted to validate `time.Time` values, `int`s, `float`s, etc?

Since we can't have generics in methods, we would have to define multiple `Check` methods, one for each type.

## Improving validation with generics

Not being able to use generics with methods is a clear limiting factor here, but there's a solution:
**just use functions instead**.

### Creating the rules package

I usually like to define a new package named `rules`, but you can name it whatever you want, even define theses functions
and interfaces in the `validator` package if you really want to

We will start by defining a generic `Rule` interface and a `RuleFunc` type.

```go
package rules

// Interface rule defines a generic validator
type Rule[T any] interface {
    Validate(context.Context, T) (bool, string)
}

// RuleFunc provides a convient way for creating Rules using functions.
type RuleFunc[T any] func(ctx context.Context, value T) (bool, string)

// Implement the Rule interface.
func (r RuleFunc[T]) Validate(ctx context.Context, value T) (bool, string) {
    return r(ctx, value)
}
```

Having it as an interface gives more flexibility in my opinion, since now a service struct can return any method
as `RuleFunc` to check if a user email is already taken on the database, for example, which will be automatically a `Rule`.

Providing `context.Context` as a parameter also allows for further customization, like adding localization, etc.

### Defining validation rules

Now that we have our interface defined, and can add some actual `Rule`s to our package:

```go
package rules

func NotBlank() RuleFunc[string] {
    return func(ctx context.Context, value string) (bool, string) {
        return validator.NotBlank(value), "Required field"
    }
}

func Email() RuleFunc[string] {
    return func(ctx context.Context, value string) (bool, string) {
        return validator.IsEmail(value), "Must be a valid email"
    }
}

func MinChars(n int) RuleFunc[string] {
    return func(ctx context.Context, value string) (bool, string) {
        return validator.MinChars(value, n), fmt.Sprintf("Must have at least %d characters.", n)
    }
}

func MaxChars(n int) RuleFunc[string] {
    return func(ctx context.Context, value string) (bool, string) {
        return validator.MaxChars(value, n), fmt.Sprintf("Must have at most %d characters.", n)
    }
}
```

Notice that we don't get rid of our `validator` helper functions. Instead, we use them for defining
our rules.

### Applying the rules

Ok, now that we defined our `Rule` interface and some implementations, it's time to define how we
are going to apply them to our validator.

To make it flexible, we will first define a `ProblemSetter` interface, which our `Validator` struct conveniently implements.

Next, we define a generic `Apply` function responsible for applying rules:

```go
package rules

// We take in any type take implements a SetProblem method,
// which is satisfied by our [validator.Validator] struct.
type ProblemSetter interface {
    SetProblem(name, reason string)
}

// Applies all rules to the provided [ProblemSetter].
func Apply[T any](ctx context.Context, ps ProblemSetter, name string, value T, rules ...Rule[T]) {
    for _, rule := range rules {
        ok, reason := rule.Validate(ctx, value)
        if !ok {
            ps.SetProblem(name, reason)
        }
    }
}
```

### The result

After defining our `rules` package, we can now update the validation in our HTTP handler. It will look like this:

```go
package main

func handleCreateUser() http.HandlerFunc {
    type request struct {
        Name            string `json:"name"`
        Email           string `json:"email"`
        Password        string `json:"password"`

        validator.Validator `json:"-"` // Omit from decoding.
    }
    return func(w http.ResponseWriter, r *http.Request) {
        var req request
        err := json.NewDecoder()
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        rules.Apply(r.Context(), &req, "name", req.Name, rules.NotBlank, rules.MaxChars(120))
        rules.Apply(r.Context(), &req, "email", req.Email, rules.NotBlank, rules.Email, rules.MaxChars(120))
        rules.Apply(r.Context(), &req, "password", req.Password, rules.NotBlank, rules.MinChars(8), rules.MaxChars(72))
        if !req.IsValid() {
            // Handle invalid request
            return
        }

        // Create new user
    }
}
```

Now, comparing the before and after

```go
// Before
req.Check(validator.NotBlank(req.Name), "name", "Required field")
req.Check(validator.MaxChars(req.Name, 120), "name", "Must have a most 120 characters")
req.Check(validator.NotBlank(req.Email), "email", "Required field")
req.Check(validator.IsEmail(req.Email), "email", "Must be a valid email")
req.Check(validator.MaxChars(req.Email, 120), "email", "Must have a most 120 characters")
req.Check(validator.NotBlank(req.Password), "password", "Required field")
req.Check(validator.MinChars(req.Password, 8), "password", "Must have at least 8 characters")
req.Check(validator.MaxChars(req.Password, 72), "password", "Must have at most 72 characters")

// After
rules.Apply(r.Context(), &req, "name", req.Name, rules.NotBlank(), rules.MaxChars(120))
rules.Apply(r.Context(), &req, "email", req.Email, rules.NotBlank(), rules.Email(), rules.MaxChars(120))
rules.Apply(r.Context(), &req, "password", req.Password, rules.NotBlank(), rules.MinChars(8), rules.MaxChars(72))
```

## Wrapping Up

By leveraging generics, we’ve made validation cleaner, more reusable, and less error-prone.
This pattern scales well across various data types and contexts.

Thanks for reading!
