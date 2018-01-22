#!/bin/bash
 
############################## PRE PROCESSING ################################
#check and process arguments
REQUIRED_NUMBER_OF_ARGUMENTS=1
if [ $# -lt $REQUIRED_NUMBER_OF_ARGUMENTS ]
then
    echo "Usage: $0 <path_to_config_file>"
    exit 1
fi

CONFIG_FILE=$1
 
echo "Config file is $CONFIG_FILE"
echo ""
 
#get the configuration parameters
source /home/ec2-user/sherlog-master/$CONFIG_FILE

############################## SETUP ################################
if [ "$USAGE" == "RUN" ]
then
	counter=1
	for node in ${VM_NODES//,/ }
	do
		echo "Running server $node ..."
		COMMAND=''		
		COMMAND=$COMMAND" fuser -k 8008/tcp;"
		COMMAND=$COMMAND" nohup go run /home/ec2-user/sherlog-master/server.go /home/ec2-user/sherlog-master/error_log $node  > /dev/null 2>&1 &"
		while [ ! -e /home/ec2-user/sherlog-master/distributed$counter.pem ] 
		do
			let counter=counter+1
		done
		ssh -p 22 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i /home/ec2-user/sherlog-master/distributed$counter.pem  -l ec2-user $node "$COMMAND"
		let counter=counter+1
	done
fi




