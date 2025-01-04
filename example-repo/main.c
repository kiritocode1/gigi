#include <stdio.h>
// map  <string, string>

#include <stdlib.h>

#define MAX_SIZE 100

typedef struct {
    int key;
    int value;
} Pair;


int main() {

    Pair *map = malloc(sizeof(Pair) * MAX_SIZE);
    int i;
    for (i = 0; i < MAX_SIZE; i++) {
        map[i].key = i;
        map[i].value = i;
    }


    printf("Hello, World!\n");

    



    return 0;
}