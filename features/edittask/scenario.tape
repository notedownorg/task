Env NOTEDOWN_DIR "features/edittask/final"
Env TEST_DATE "2020-01-01"

Source "features/common.tape"
Output "features/edittask/demo.gif"

Sleep 2

# Complete a recurring task
Type "e"
Sleep 2
Type "x"
Sleep 2
Enter
Sleep 2

# Show the recurrence
Right 
Sleep 2
Left
Sleep 1

# Mark the blocked task as in progress
Down 2
Type "e"
Type "/"
Tab
Right@1ms 43
Backspace@1ms 43
Type "This task is no longer blocked!"
Sleep 1
Enter

# Reschedule the unblocked task to tomorrow
Up
Sleep 1
Type "r"
Sleep 1
Type "1"
Sleep 1
Right
Sleep 2
Left

Sleep 5
