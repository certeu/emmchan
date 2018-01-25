# emmchan #

emmchan is a command-line utility which adds new channels to an EMM channel
directory. emmchan loads an existing channel directory, reads new feed URL
from STDIN and adds it to the directory. On exit the new channel directory is
written to STDOUT.

## Usage ##

### Bulk ###

```sh
echo "https://rss-feed-url" | emmchan -d channeldirectory.xml > out.xml
cat n.txt | emmchan -d channeldirectory.xml > out.xml
```

```cmd
echo http://rss-feed-url| emmchan.exe -d channeldirectory.xml > out.xml
type n.txt | emmchan.exe -d channeldirectory.xml > out.xml
```
### Interactive ###

```sh
$ ./emmchan -d channeldirectory.xml
$ https://rss-feed-url
$ ^D
```

```cmd
$ emmchan.exe -d channeldirectory.xml
$ https://rss-feed-url
$ ^Z
```
