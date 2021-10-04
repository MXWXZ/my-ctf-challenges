#include "BP.h"
#include <algorithm>

struct Res {
    char chr;
    double value;
};

int main() {
    srand((unsigned)time(0));
#ifndef TEST
    FILE* fp = _tfopen(_T("./data.txt"), _T("r"));
    Vector<Data> datga;
    char ch;
    while ((ch = fgetc(fp)) != EOF) {
        Data t;
        for (int i = 0; i < 26; ++i)
            if (i != ch - 'a')
                t.y[i] = 0;
            else
                t.y[i] = 1;
        for (int i = 0; i < 300; ++i)
            t.x[i] = fgetc(fp) - '0';
        data.push_back(t);
        fgetc(fp);  // eat /n
    }

    fclose(fp);
#endif

    BP* bp = new BP();
#ifndef TEST
    bp->GetData(data);
#ifdef GOON
    bp->GoOn();
#else
    bp->Train();
#endif
#endif

#ifdef TEST
    bp->Load();
    char ch;
    while ((ch = fgetc(stdin)) != '2') {
        Vector<bool> in;
        in.push_back(ch - '0');
        for (int i = 1; i < 300; ++i)
            in.push_back(fgetc(stdin) - '0');
        Vector<double> out;
        out = bp->ForeCast(in);
        Res res[26];
        for (int i = 0; i < 26; ++i) {
            res[i].chr = 'a' + i;
            res[i].value = out[i];
        }
        std::sort(res, res + 26, [](const Res& r1, const Res& r2) {
            return r1.value > r2.value;
        });
        printf("%c", res[0].chr);
    }
#endif
#ifndef TEST
    getchar();
#endif
    return 0;
}
