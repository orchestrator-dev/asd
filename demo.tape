Output assets/demo.gif

Require ./dist/asd

Set FontSize 15
Set Width 1200
Set Height 800
Set Padding 20
Set Theme "Dracula"
Set Framerate 24

Hide
Type "clear"
Enter
Show

Type "# Let's explore the power of asd!"
Sleep 1s
Enter
Enter

Type "# 1. Beautiful Syntax Highlighting & Line Numbers"
Enter
Type "./dist/asd cmd/root.go | head -n 25"
Sleep 1s
Enter
Sleep 3s

Type "clear"
Enter

Type "# 2. Intelligent Log Parsing and Colorization"
Enter
Type "./dist/asd test.log"
Sleep 1s
Enter
Sleep 3s

Type "clear"
Enter

Type "# 3. Native Jupyter Notebook Rendering"
Enter
Type "./dist/asd test.ipynb"
Sleep 1s
Enter
Sleep 3s

Type "clear"
Enter

Type "# 4. Instant SQLite Database Exploration"
Enter
Type "./dist/asd test.db"
Sleep 1s
Enter
Sleep 3s

Type "clear"
Enter

Type "# 5. Directory Dashboard with Git Integration"
Enter
Type "./dist/asd ."
Sleep 1s
Enter
Sleep 5s
