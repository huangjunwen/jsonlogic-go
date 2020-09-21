## jl

Command line json logic evaluator.

Pass logic and (optional) data to the stdin to see ouput in stdout.

Example usage:

```bash
$ echo '{"+":[{"var":"left"}, {"var":"right"}]} {"left":3, "right":4}' | jl
7
$ echo '{"var":""}' | jl
{}
```
