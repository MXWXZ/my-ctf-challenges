perl -p preprocess.pl < map.h > map.ph
i686-w64-mingw32-gcc -D LANG_EN -Wall -O3 -o tsw_en.exe tsw.c icon.dll -mwindows -lgdi32 -lwinmm
i686-w64-mingw32-gcc -D LANG_CH -Wall -O3 -o tsw_ch.exe tsw.c icon.dll -mwindows -lgdi32 -lwinmm
rm -rf map.ph
mv tsw_ch.exe ../
mv tsw_en.exe ../