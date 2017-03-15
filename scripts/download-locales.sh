#!/usr/bin/bash
cd ../data/locales
rm templates-*.po
for i in it es de pl ru ko en_GB; do
	echo "$i"
	wget --content-disposition --quiet "https://cutebirbs.ripple.moe/export/?path=/$i/Hanayo/"
done
