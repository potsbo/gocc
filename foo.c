#include <stdio.h>
#include <stdlib.h>

int foo() { return 0; }
int bar(int v) { return v; }
int add(int a, int b) { return a + b; }
int print(int a) { return printf("%#010x\n", a); }
void alloc(int **p, int a, int b, int c, int d) {
	int *pt = (int *) malloc(4 * sizeof(int));
	pt[0] = a;
	pt[1] = b;
	pt[2] = c;
	pt[3] = d;
	*p = pt;
}
