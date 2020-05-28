set -o errexit
set -o pipefail

# TODO change to your desired driver.
## VENDOR=${VENDOR:-k8s}
## DRIVER=${DRIVER:-dummy}

VENDOR=k8s
DRIVER=nfs

# Assuming the single driver file is located at /$DRIVER inside the DaemonSet image.

driver_dir=$VENDOR${VENDOR:+"~"}${DRIVER}
if [ ! -d "/flexmnt/$driver_dir" ]; then
  mkdir "/flexmnt/$driver_dir"
fi

tmp_driver=.tmp_$DRIVER
cp "/$DRIVER" "/flexmnt/$driver_dir/$tmp_driver"
mv -f "/flexmnt/$driver_dir/$tmp_driver" "/flexmnt/$driver_dir/$DRIVER"

while : ; do
  sleep 3600
done
