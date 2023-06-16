#include <stdio.h>

char return_full_name(char first[30], char last[30])
{
    char full_name[60];
    int i, j;

    for (i = 0; first[i] != '\0'; i++)
    {
        full_name[i] = first[i];
    }
    full_name[i] = ' ';
    i++;
    for (j = 0; last[j] != '\0'; j++)
    {
        full_name[i + j] = last[j];
    }
    full_name[i + j] = '\0';
    return full_name;
}

int main()
{
    char first_name[] = "John";
    char last_name[] = "Doe";

    char full_name = return_full_name(first_name, last_name);
    printf("Hello, %s!\n", full_name);
    return 0;
}