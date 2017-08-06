# go_doorbell
Reimplementation of the doorbell but in Go.


I'm redoing the doorbell.
This time it's in Go. My first foray into the language so not a magnificent exemplar of good Golang.

It will use wifi to send ring commnds to each wireless speaker.

At present I still need to write some scripts to do some of this stuff.
And something that will auto run them etc

## To Set Up Central
* Make the dir into which you will run things from now we'll call it $THEDIR
* make a dir called "thesounds" off that
* Copy the "doorbell_button" and "doorbell_manage" executables to $THEDIR
* Create a file called "config.json" And put in ```{
"Doorbell_dir":"$THEDIR",
"Satellite_port":3400
}```
* ```touch doorbells.txt``` To create the doorbells list.
* run "./doorbell_button"
* run "./doorbell_manage"
* to set things up for autorunning on reboot you will need to edit the "rc.local" file in ```/etc```. Remember to cd to the dir before running. And you may want to redirect output to a log file with the double arrow(>>) to append. Bash scripts to do these are in the repo.
* To upload sound files point your browser at http://Server_name_or_IP:3434. There is a link to the upload form there. If there are already subscribed doorbells, the files should be pushed out to them all.


## To Set Up The Satellites
* Make the dir into which you will run things from now we'll call it $THEDIR
* make a dir called "thesounds" off that
* Copy the "doorbell" executable to $THEDIR. If this is a ARM based CPU like a Raspberry Pi then you will have copy the executable called "armdoorbell"
* Copy the "sync_chimes" executable to $THEDIR. If this is a X86 machine you will need the "xsynch_chimes" file.
* Create a file called "config.json" And put in ```{
"Doorbell_dir":"$THEDIR",
"Satellite_port":3400
}```
* run ".doorbell" (or "./armdoorbell")
* open another login window to doorbell machine
* go to $THEDIR
* run ./sync_chimes $THEIPOFTHECENTRAL (Please note that the doorbell executable must be running while this is run)
* to set things up for autorunning on reboot you will need to edit the "rc.local" file in ```/etc```. Remember to cd to the dir befor running. And yoy may want to redirect output to a log file with the double arrow(>>) to append.  






