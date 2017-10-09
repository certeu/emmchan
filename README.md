# emmparser #

emmparser is a command-line utility which adds new channels to an EMM channel
directory. emmparser loads an existing channel directory, reads new feed URL
from STDIN and adds it to the directory. On exit the new channel directory is
written to STDOUT.

## Usage ##

```sh
echo "https://rss-feed-url" | emmparser -d channeldirectory.xml > out.xml
cat n.txt | emmparser -d channeldirectory.xml > out.xml
```

```cmd
echo http://rss-feed-url| emmparser.exe -d channeldirectory.xml > out.xml
type n.txt | emmparser.exe -d channeldirectory.xml > out.xml
```

```sh
    $ ./emmparser -d channeldirectory.xml
    $ https://rss-feed-url
    $ ^D
```

```cmd
    $ emmparser.exe -d channeldirectory.xml
    $ https://rss-feed-url
    $ ^Z
```
