# Service entry point

`main.go` will serve as the entry point for starting both REST and gRPC adapters

## Considerations for this folder

- Here's where main applications for this project reside.
- Don't add a lot of code in this folder. If code can be reused/imported
  elsewhere then it should live in `/pkg`.
- If you need the code but you don't want others to reuse it, then
  the code should live in `/internal`.

As a rule of thumb:

- Keep a small `main` function.
- Import code from `/internal` and/or `/pkg`
