# Advent of code 2024

```sh
# Run solution to day 1, sub-task 1, etc.
go run . 1_1
# Type or copy-paste the input right into the terminal

# Or better yet do this:
cat input.txt | go run . 1_1
```

Day 6 can be run step by step. Create a `patrol_map.txt` file in the current
working directory and run `go run . 6_1` (or `./main 6_1`) repeatedly. You can
use `watch` to have a somewhat animated experience like so:
```sh
watch -n .1 ./main 6_1
```
Unfortunately `watch` doesn't allow you to scroll the program's output. If only
there was a program for that, [cough cough](https://github.com/gmoshkin/dotfiles/raw/refs/heads/master/jai/pageview)
