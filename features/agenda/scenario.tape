Env NOTEDOWN_DIR "features/agenda/final"
Env TEST_DATE "2020-01-01"

Source "features/common.tape"
Output "features/agenda/demo.gif"

Sleep 2

# Complete the in progress task
Type "x"
Sleep 2

# Reschedule the todo task to tomorrow
Type "r"
Sleep 1
Type "1"
Sleep 1
Right
Sleep 2
Left
Sleep 2

# Delete the blocked task
Type "d"

Sleep 5
