This program extracts data from the Covid tracking project and does one of two
things with the data:
1. Writes a csv file so it can be manipulated by a spreadsheet or similar application
2. Plot the data on the screen using a browser

The first case is handled in the covid/cmd/core folder.  The second case is
handled in the web folder.  Both of these programs depend on covid/pkg/virusdata
package which extracts and parses the returned JSON files.  It in turn depends
on Lexer2 from github.com/Saied74/Lexer2.

The programs are pretty simple and well commented (I think), so I won't belabor
the points anymore.

To do list:
1. At least put the filenames in the constants folder or a config file
2. Handle no data return and other exceptions in the web side the web way
