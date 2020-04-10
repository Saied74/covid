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

The application is configured using a file called config.csv.  You can change
the name of this file with the -c flag.  An example of this file is shown at
the end of this readme files.  The program searches for this file in one of the
following ways and in this exact order:
1. It looks in the same directory where the application is starting
2. It looks in the $HOME/src/covid directory
3. It looks in a directory pointed to by the environment variable of -c flag

The environment variable given with -c flag can be upper or lower case, with
or without $ sign.


Sample config.csv
```

# Be careful with typing.  This file controls lots of things
# All comments look like this
# The first two columns are used by the program, the rest are ignored.
# The first column has to be exactly as shown
# Symbol | is used instead of commas so they can be used without escaping
# Blank lines are ignored
# All addresses are relative to appHome
# appHome can start with environment variable or be absolute
# Environment variable must start with $ sing
# The "/" joining addresses is provided by the program
# Lines can appear in any order, they are processed together

appHome|$HOME/Documents/gocode/src/covid|
patternFile|config/pattern.csv| used to process the JSON data by the lexer
csvOutputFile|config/outreq.csv| used to select and format the csv output
covidProjectURL|https://covidtracking.com/api/states/daily

# there can be as many template files as needed
templateFile|ui/html/base.page.tmpl
templateFile|ui/html/plot.partial.tmpl

# one of the templateFiles must also be a plotFile to be used for plotting
plotFile|ui/html/plot.partial.tmpl

# server ip address
ipAddress|:8080
```
