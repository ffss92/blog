---
title: Leveraging Go generics for input validation
subtitle: |
    Go's generics may not be as extensive as those in other languages, but they still offer powerful capabilities. 
    In this article, I'll demonstrate how I utilize generics to simplify and enhance validation in my projects.
author: "@ffss"
draft: false
date: "2024-12-01"
tags:
    - Go
    - Generics
---

## A simple, but effective, pattern.

[Alex Edwards](https://www.alexedwards.net/) introduces us in his excellent [Let's Go](https://lets-go.alexedwards.net/)
books a simple way of validating HTTP requests in Go. It starts by defining a `validator` struct that looks something like this:

```go
type validator struct {
    problems map[string]string
}

func (v *validator) setProblem(name, reason string) {
    if v.problems == nil {
        v.problems = make(map[string]string)
    }
    if _, ok := v.problems[name]; !ok {
        v.problems[name] = reason
    }
}

func (v *validator) check(ok bool, name, reason string) {
    if !ok {
        v.setProblem(name, reason)
    }
}

func (v validator) isValid() bool {
    return len(v.problems) == 0
}
```

Very straight fowarded, no tag magic, just a simple struct that collect validation
problems and sets them in a map. And, honestly, it's good enough for a lot of cases.

To meet the `ok` check in the validator struct, you would define some helper functions,
like this:

```go
func notBlank(value string) bool {
    // Check if string if not empty
}

func isEmail(value string) bool {
    // Check if its an actual email
}

func maxLength(value string, n int) bool {
    // Check if string has at most n characters
}

// Etc...
```

Finally, putting it all together, a HTTP request validation would look something like this:

```go
package handler 

type routeRequest struct { 
    Name string `json:"name"`
    Email string `json:"email"`
    validator // embed validator struct, if it is public set `json:"-"`
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
    var req routeRequest
    err := decode(r, &req)
    if err != nil {
        // Handle error
    }

    req.check(notBlank(req.Name), "name", "Required field")
    req.check(maxLength(req.Name, 255), "name", "Must have at most 255 characters")
    req.check(notBlank(req.Email), "email", "Required field")
    req.check(isEmail(req.Email), "email", "Must be a valid email")
    if !req.isValid() {
        // Handle invalid payload
    }
    // Otherwise, continue
}
```

It works fine for a starting point, but wouldn't it be nice if we could:

- Standardize messages
- Perform multiple checks in a single line
- Add some context, like preferred language, etc.

A naive first approach to solving this would be pass the helper functions as variadic
parameters, something like this:

```go
req.check("email", req.Email, notBlank, isEmail)
```

Except now we are no longer informing the validator what the error message is. We can fix this
by changing the checker function signature to something like this:

```go
// In helper functions
func isEmail(value string) (bool, string) {
    // Check if it's an email
    return ok, "Must be a valid email"
}

// Our maxLength function has to meet the signature now...
func maxLength(n int) func(string) (bool, string) {
    return func(value string) {
        // Check length
        return ok, fmt.Sprintf("Must have at most %d characters", n)
    }
}

// In validator
func (v *validator) check(name, value string, checkers ...func(string) (bool, string)) {
    for _, checker := range checkers {
        ok, reason := chck(value)
        if !ok {
            v.setProblem(name, reason)
        }
    }
}
```

This works great! Except it works for `string`'s only... What if we need to validate an `int` value, and
an `int32` and `int64` maybe? Don't forget about all of those `uint`'s and `float`'s too. 

Normally, we would just slap a generic on this method, but since we can't have generics in Go struct methods,
we would have to do something like this:

```go
func (v *validator) checkString(name, value string, checkers ...func(string) (bool, string)) {
    for _, checker := range checkers {
        ok, reason := checker(value)
        if !ok {
            v.setProblem(name, reason)
        }
    }
}

func (v *validator) checkInt(name, value int, checkers ...func(int) (bool, string)) {
    for _, checker := range checkers {
        ok, reason := checker(value)
        if !ok {
            v.setProblem(name, reason)
        }
    }
}

func (v *validator) checkTime(name, value time.Time, checkers ...func(time.Time) (bool, string)) {
    for _, checker := range checkers {
        ok, reason := checker(value)
        if !ok {
            v.setProblem(name, reason)
        }
    }
}

// Ad infinitum...
```

## A different, generic approach

After spending some time playing with the idea of using generics for validation library, I came to the
conclusion that it just wasn't worth it. 

Instead, I came up with some sort of small framework for creating validation libraries, and it's actually
quite simple, as it ended up being a very small API, it consists of defining rules and applying them.

I started by defining a `Rule` interface:

```go
// Package rules define a framework for building validation packages.
package rules

type Rule[T any] interface {
    // Apply returns whether value is valid, and the reason if it is not.
    Apply(ctx context.Context, value T) (bool, string)
}
```

In the interface above, `T` is a generic type, so it could be pretty much anything. To make it easier to
create rules, I also created a `RuleFunc` type that implements the `Rule` interface:

```go
package rules

// Check at compile time if RuleFunc implements Rule.
var _ Rule[any] = (*RuleFunc[any])(nil)

type RuleFunc[T any] func(ctx context.Context, value T) (bool, string)

func (r RuleFunc[T]) Apply(ctx context.Context, value T) (bool, string) {
    return r(ctx, value)
}
```

This means we could easily create a rule with your previously defined funcs like this: 

```go
// Package comp defines funcions for cheking if a value meet certain criteria.
package comp

func NotBlank(value string) bool {
    return utf8.RuneCountInString(value) > 0
}
```

```go
// Package validator defines rules for validation.
package validator

func NotBlank() RuleFunc[string] {
    return func(ctx context.Context, value string) (bool, string) {
        comp.NotBlank(value), "Field is required"
    }
}
```

You could also use the context to generate different reason message, etc. Here's a very simple example:

```go
package validator 

func NotBlank() RuleFunc[string] {
    return func(ctx context.Context, value string) (bool, string) {
        var reason string
        switch language.FromContext(ctx) {
        case "pt-BR":
            reason = "Campo obrigatório"
        // ...
        default:
            reason = "Field is required"
        }

        return comp.NotBlank(value), reason
    }
}
```

Finally, we define an `Apply` function that will perform the application
of rules in a `ProblemSetter`:

```go
package rules

type ProblemSetter interface {
    SetProblem(name, reason string)
}

// Applies all rules.
func Apply[T any](ctx context.Context, setter ProblemSetter, name string, value T, rules ...Rule[T]) {
	for _, rule := range rules {
		ok, reason := rule.Apply(ctx, value)
		if !ok {
			setter.SetProblem(name, reason)
			return
		}
	}
}

// Applies rules conditionally
func ApplyIf[T any](ctx context.Context, setter ProblemSetter, cond bool, name string, value T, rules ...Rule[T]) {
	if cond {
	    Apply(ctx, setter, name, value, rules...)
	}
}
```

## Putting it all together

Using our new approach we would have the perform some small changes in our original code. We'll start
by implementing `ProblemSetter`. Lets define a new `validator` package. You could do this anyway you
want, but here's an example:

```go
package validator

var _ rules.ProblemSetter = (*validator)(nil)

type Validator struct {
    problems map[string]string
}

// This is easy. It used to be called `setProblem`, just capitalize the first 's'.
func (v *Validator) SetProblem(name, reason string) {
    if v.problems == nil {
        v.problems = make(map[string]string)
    }
    if _, ok := v.problems[name]; !ok {
        v.problems[name] = reason
    }
}

func (v *Validator) IsValid() bool {
    return len(v.problems) == 0
}
```

Then, we create our rules, also in the validator package:

```go
package validator

func NotBlank() rules.RuleFunc[string] {
	return func(ctx context.Context, value string) (bool, string) {
		return cond.NotBlank(string(value)), "Field is required."
	}
}

func Email() rules.RuleFunc[string] {
	return func(ctx context.Context, value string) (bool, string) {
		return cond.IsEmail(value), "Must be a valid email."
	}
}

func MaxChars(n int) rules.RuleFunc[string] {
	return func(ctx context.Context, value string) (bool, string) {
		return cond.MaxChars(value, n), fmt.Sprintf("Must have at most %d characters.", n)
	}
}

func GreaterThan[T cmp.Ordered](target T) rules.RuleFunc[T] {
	return func(ctx context.Context, value T) (bool, string) {
		return cond.GreaterThan(value, target), fmt.Sprintf("Must be greater than %v.", target)
	}
}
```

Now, we can bring this all together. Out validation process will look like this:

```go
package handler 

type routeRequest struct { 
    Name string `json:"name"`
    Email string `json:"email"`
    Age uint8 `json:"age"`

    validator.Validator `json:"-"`
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
    var req routeRequest
    err := decode(r, &req)
    if err != nil {
        // Handle error
    }

    rules.Apply(r.Context(), &req, "name", req.Name, validator.NotBlank(), validator.MaxLength(255))
    rules.Apply(r.Context(), &req, "email", req.Email, validator.NotBlank(), validator.Email(), validator.MaxLength(255))
    rules.Apply(r.Context(), &req, "age", req.Age, validator.GreaterThan(17))
    if !req.IsValid() {
        // Handle invalid payload
    }
    // Otherwise, continue
}
```

From this:

```go
package handler 

type routeRequest struct { 
    Name string `json:"name"`
    Email string `json:"email"`
    Age uint8 `json:"age"`

    validator // embed validator struct, if it is public set `json:"-"`
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
    var req routeRequest
    err := decode(r, &req)
    if err != nil {
        // Handle error
    }

    req.check(notBlank(req.Name), "name", "Required field")
    req.check(maxLength(req.Name, 255), "name", "Must have at most 255 characters")
    req.check(notBlank(req.Email), "email", "Required field")
    req.check(isEmail(req.Email), "email", "Must be a valid email")
    req.check(maxLength(req.Email, 255), "email", "Must have at most 255 characters")
    req.check(req.Age > 17, "age", "Must be greater than 17")
    if !req.isValid() {
        // Handle invalid payload
    }
    // Otherwise, continue
}
```

This is it! Thanks for reading.
