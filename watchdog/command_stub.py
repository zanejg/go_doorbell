import sys

if __name__ == "__main__":
    print("there were {} parms and they were:{}".format(len(sys.argv),sys.argv))


    if len(sys.argv) == 2:
        doorbell_filename = sys.argv[1]
    else:
        sys.exit("Usage: watchdog.py filename")
        

    # first get a list of the subscribed doorbell clients
    try:
        with open(doorbell_filename,"r") as doorbell_file:
            # we also want to exclude the local machine or empty lines
            doorbells = [l.strip() for l in doorbell_file
                        if l.strip() not in ["127.0.0.1","localhost",''] ]
    except IOError:
        sys.exit('Doorbell list file "{}" does not exist'.format(doorbell_filename))
    
    

