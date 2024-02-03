#!/usr/bin/python
#############################################
# To periodically poll the subscribed doorbells and either 
# light up the LED or turn it off
###########################################################
import os
from gpiozero import LED
import sys
import time
import requests

doorbells = []
led = LED(4)

print("start init delay")

# on first start up we need to ensure that things have settled down for
# our watch dog to start so we will wait for 60 seces before we start
time.sleep(15)

print("ended init delay")

if __name__ == "__main__":
    print("name = main")
    if len(sys.argv) == 2:
        doorbell_filename = sys.argv[1]
    else:
        sys.exit("Usage: watchdog.py filename")
        
    print("opened doorbells file")

    # first get a list of the subscribed doorbell clients
    try:
        with open(doorbell_filename, "r") as doorbell_file:
            # we also want to exclude the local machine or empty lines
            doorbells = [l.strip() for l in doorbell_file
                        if l.strip() not in ["127.0.0.1","localhost",''] ]
    except IOError:
        sys.exit('Doorbell list file "{}" does not exist'.format(doorbell_filename))
        
    print("created doorbells array")

    #print(doorbells)

    # now we have a list of subscribed doorbells from the original source 
    # we can check that we can see them 

    while True:
        # we need to check that the sever process is up and healthy too
        # so we will send a get request to it
        print("about to check localhost")
        resp = requests.get("http://127.0.0.1:3434")
        if not (resp.status_code == 200 and 
                "Doorbell Central Controller" in resp.text):
            # the server process is not working
            # so we really need to reboot
            # but this could mean we might get into a nasty reboot loop so 
            # we will put a delay to allow us to kill it 
            print("Server will reboot in 30 secs!!")
            time.sleep(30)
            import os
            os.system('sudo shutdown -r now')
        
        print("Doorbell server process appears to be OK")
    

        # at first I will only check that I can see one other as I am really 
        # just checking that the network is up
        # so start with OK being false
        ok = False
        for this_doorbell in doorbells:
            response = os.system("ping -c 1 {} > /dev/null".format(this_doorbell))
            connected = (response == 0)

            if(connected):
                ok=True
                print("{} is up".format(this_doorbell))
            else:
                print("{} is down".format(this_doorbell))

        if ok:
            led.on()
            print("####### OK  ######")
        else:
            led.off()
            print("!!!! NOT OK !!!!")
        
        time.sleep(60)
