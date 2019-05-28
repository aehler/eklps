#!/bin/bash

PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/opt/eklps/bin/

echo "getting horses 1-12000"

for i in {5845..5848}
do
   curl -d "racer=$i" --cookie 'PHPSESSID=n7ng96d40lepb65g4a7k50dpb7'  http://eklps.com/stables/info_horsecard.php >> ./horsedata/$i.horse 
done

exit 0
