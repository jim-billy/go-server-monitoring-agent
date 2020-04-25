#!/bin/bash
#set -x
SUCCESS=0
FAILURE=1
LINUX_DISTRIBUTION=''
OS_NAME=`uname -s`

export LC_NUMERIC="en_US"

setOperatingSystemAndDistribution() {
	if [ "${OS_NAME}" = "SunOS" ]; then
		OS_NAME=Solaris		
	elif [ "${OS_NAME}" = "AIX" ]; then
		OS_NAME=AIX	
	elif [ "${OS_NAME}" = "Linux" ]; then				
		if [ -f /etc/redhat-release ]; then
			LINUX_DISTRIBUTION="RedHat"
		elif [ -f /etc/centos-release ]; then
			LINUX_DISTRIBUTION="CentOS"		
		elif [ -f /etc/SuSE-release ]; then
			LINUX_DISTRIBUTION="Suse"		
		elif [ -f /etc/mandrake-release ]; then
			LINUX_DISTRIBUTION="Mandrake"
		elif [ -f /etc/fedora-release ]; then
			LINUX_DISTRIBUTION="Fedora"			
		elif [ -f /etc/debian_version ]; then
			LINUX_DISTRIBUTION="Debian"			
		fi
	fi

}


cpuUtil() {
    topOutput=`top -b -d1 -n2 | awk 'BEGIN{ORS=" ::: ";} /^%Cpu|^Cpu/ '`
    #echo $topOutput        
    echo $topOutput | awk 'BEGIN{FS=" ::: ";}
    {
        for(k=(NF/2+1);k<=NF;k++)
        {
            split($k,cpuIdleArr,"wa")
           	split(cpuIdleArr[1],niceArr,"ni")
			split(cpuIdleArr[1],waitArr,"id")
			awkCpu = substr(niceArr[2],2)
			awkWait = substr(waitArr[2],2)
            gsub(",",".",awkCpu)
			gsub(",",".",awkWait)
            totIdlePerc+=awkCpu
        }
        cpu_idle_percentage = totIdlePerc/(NF/2)
        cpu_util = 100.0 - cpu_idle_percentage
        print "cpu_instance : cpu"
        print "cpu_idle_percentage : ",cpu_idle_percentage           
        printf "cpu_load_percentage : %.2f\n",cpu_util
		printf "cpu_wait_percentage : %.2f\n",awkWait
    }'
               
}

memUtil() {
	free -k | awk 'BEGIN {  
        ORS="\n";
	virtTrue=0;
	avail = 0;
	}
	{
	mem = match($1,"Mem")
	swap = match($1,"Swap")
	if(NR == 1)
	{
		avail = match($0,"available")
		#print "Setting available to :",avail	
	}
        if (mem == 1)
            { 
		print "total_visible_memory :",$2
		if( avail != 0 )
		{
			print "free_physical_memory :"$7	
		}
		else
		{
			print "free_physical_memory :",($2-$3+$6+$7)
		}
	    }	
        else if (swap == 1)
             {print "total_virtual_memory :",$2,"\nfree_Virtual_memory :",$4
	    virtTrue = 1
	    }       
	}END{if(virtTrue == 0) print "total_virtual_memory :",0,"\nfree_Virtual_memory :",0}'
	setOperatingSystemAndDistribution
	strcmd="/bin/cat /etc/issue"
	if [ ${LINUX_DISTRIBUTION} = "RedHat" ]; then
		strcmd="/bin/cat /etc/redhat-release"
	elif [ ${LINUX_DISTRIBUTION} = "CentOS" ]; then
		strcmd="/bin/cat /etc/centos-release"
	else 
		strcmd="/bin/cat /etc/issue"	
	fi

	$strcmd  | awk 'BEGIN {  
        ORS="\n";
        LINUX_DISTRIBUTION="";
	}
	{
        if (NR<6 && $0 != "") {
            LINUX_DISTRIBUTION=LINUX_DISTRIBUTION""$0
            gsub(/\\n/,"",LINUX_DISTRIBUTION)
            gsub(/\\l/,"",LINUX_DISTRIBUTION)
            gsub(/\\r/,"",LINUX_DISTRIBUTION)
            gsub(/\\m/,"",LINUX_DISTRIBUTION)	                
        }       
	}END {print "os_name : ",LINUX_DISTRIBUTION}'
}

diskUtil() {
    df -l -T | grep -ivE 'Filesystem|overlay' | awk 'BEGIN {
        ORS="\n";
        partitionNameArray[0] = 0;
    }
    {       
        #print "Processing Record ",NR,NF,$NF;
        if ($NF in partitionNameArray == 0)
        {
            partitionNameArray[$NF] = $NF
            if (NF > 2)
            {
				printf $NF" == %.0f == %.0f \n",(($(NF-3)*1024)+($(NF-2)*1024)), $(NF-2)*1024      
            }
        }
    }'
}

parseInputParam() {
	if [ "$1" != "" ]; then
		for INPUT_PARAM in ${1//,/ }; do
			if [ "${INPUT_PARAM}" = "cpu_util" ]; then
				cpuUtil
			elif [ "${INPUT_PARAM}" = "mem_util" ]; then
				memUtil
			elif [ "${INPUT_PARAM}" = "disk_util" ]; then
				diskUtil
			fi
		done
	fi
}

main() {
	parseInputParam $1
}

main $1
