# go_doorbell
Reimplementation of the doorbell but in Go.


I'm redoing the doorbell.
This time it's in Go. My first foray into the language so not a magnificent exemplar of good Golang.

It will use wifi to send ring commnds to each wireless speaker.

At present I still need to write some scripts to do some of this stuff

## To Set Up Central
* Make the dir into which you will run things from now we'll call it $THEDIR
* make a dir called "thesounds" off that
* Copy the "doorbell_central" executable to $THEDIR
* Create a file called "config.json" And put in ```{
"Doorbell_dir":"$THEDIR",
"Satellite_port":3400
}```
* ```touch doorbells.txt``` To create the doorbells list.
* run "./doorbell_central"






