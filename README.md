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

I have gotten rid of the complex set up and verification code and replaced it
with an environment variable and two new flags.  The old flags are gone.

The environment variable is PATHVID and it must be set to the head of the
project which in my case is $GOPATH/src/covid.  In the Dockerfile it is set
to /go/covid.

The two flags are covidProjectURL and ipAddress.  covidProjectURL overrides
the URL for getting the covid porject data.  The ipAddress overrides the server
IP address that the web server listens to.
