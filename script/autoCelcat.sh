pwd=`pwd`
celcat="${pwd}/../CY-celcat"

if ! [ -d "./config" ]; then
    mkdir config
fi

if ! [ -d "./out" ]; then
    mkdir out

fi

cd out

for f in ../config/*.json; do
    f_name=${f##*/}
    f_name=${f_name%.*}
    echo $f_name
    $celcat -c $f -svg "${f_name}.svg"
    inkscape -w 1920 -h 1080 "${f_name}.svg" -e "${f_name}.png"
    $celcat -c $f
    mv data.ics "${f_name}.ics"
done

cd ..

cp out/* /home/ubuntu/docker/main/web-nginx/
