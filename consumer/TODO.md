### timed stats
- align to wallclock
- better way to store arbitrary stats (ex:
  https://github.com/mindreframer/golang-stuff/tree/master/github.com/tumblr/gocircuit/src/circuit/kit/stat
  )

### aggregated stats
- reimplement rbtree, balancing in-place

### alerting
- should emit alerts to some other entity to take care

### general
- count stats towards log time, not towards `time.Now()`
- use termbox to look more htop-ish
- expose stats to lua/otto to create complex alert rules
