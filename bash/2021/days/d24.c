#include <stdio.h>
#include <stdlib.h>
long int A, B, C, D, E, F, G, H, I, J, K ,L, M, N;
long int Z0=0, Z1, Z2, Z3, Z4, Z5, Z6, Z7, Z8, Z9, Z10, Z11, Z12, Z13, Z14;

int exercise1(){
    printf("Exercise 1:\n");
    for(A=9;A>0;A--){
        Z1=Z0*26+A;
        for(B=9;B>0;B--){
            Z2=Z1*26+B+6;
            for(C=9;C>0;C--){
                Z3=Z2*26+C+4;
                for(D=9;D>0;D--){
                    Z4=Z3*26+D+2;
                    for(E=9;E>0;E--){
                        Z5=Z4*26+E+9;
                        for(F=9;F>0;F--){
                            Z6=(Z5/26*((25*(Z5%26+-2==F==0))+1))+((F+1)*(Z5%26+-2==F==0));
                            for(G=9;G>0;G--){
                                Z7=Z6*26+G+10;
                                for(H=9;H>0;H--){
                                    Z8=(Z7/26*((25*(Z7%26+-15==H==0))+1))+((H+6)*(Z7%26+-15==H==0));
                                    for(I=9;I>0;I--){
                                        Z9=(Z8/26*((25*(Z8%26+-10==I==0))+1))+((I+4)*(Z8%26+-10==I==0));
                                        for(J=9;J>0;J--){
                                            Z10=Z9*26+J+6;
                                            for(K=9;K>0;K--){
                                                Z11=(Z10/26*((25*(Z10%26+-10==K==0))+1))+((K+3)*(Z10%26+-10==K==0));
                                                for(L=9;L>0;L--){
                                                    Z12=(Z11/26*((25*(Z11%26+-4==L==0))+1))+((L+9)*(Z11%26+-4==L==0));
                                                    for(M=9;M>0;M--){
                                                        Z13=(Z12/26*((25*(Z12%26+-1==M==0))+1))+((M+15)*(Z12%26+-1==M==0));
                                                        for(N=9;N>0;N--){
                                                            Z14=(Z13/26*((25*(Z13%26+-1==N==0))+1))+((N+5)*(Z13%26+-1==N==0));
                                                            if (Z14==0) {
                                                                printf("%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld\n", A, B, C, D, E, F, G, H, I, J, K ,L, M, N);
                                                                exit(0);
                                                            }
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

int exercise2(){
    printf("Exercise 2:\n");
    for(A=1;A<10;A++){
        Z1=Z0*26+A;
        for(B=1;B<10;B++){
            Z2=Z1*26+B+6;
            for(C=1;C<10;C++){
                Z3=Z2*26+C+4;
                for(D=1;D<10;D++){
                    Z4=Z3*26+D+2;
                    for(E=1;E<10;E++){
                        Z5=Z4*26+E+9;
                        for(F=1;F<10;F++){
                            Z6=(Z5/26*((25*(Z5%26+-2==F==0))+1))+((F+1)*(Z5%26+-2==F==0));
                            for(G=1;G<10;G++){
                                Z7=Z6*26+G+10;
                                for(H=1;H<10;H++){
                                    Z8=(Z7/26*((25*(Z7%26+-15==H==0))+1))+((H+6)*(Z7%26+-15==H==0));
                                    for(I=1;I<10;I++){
                                        Z9=(Z8/26*((25*(Z8%26+-10==I==0))+1))+((I+4)*(Z8%26+-10==I==0));
                                        for(J=1;J<10;J++){
                                            Z10=Z9*26+J+6;
                                            for(K=1;K<10;K++){
                                                Z11=(Z10/26*((25*(Z10%26+-10==K==0))+1))+((K+3)*(Z10%26+-10==K==0));
                                                for(L=1;L<10;L++){
                                                    Z12=(Z11/26*((25*(Z11%26+-4==L==0))+1))+((L+9)*(Z11%26+-4==L==0));
                                                    for(M=1;M<10;M++){
                                                        Z13=(Z12/26*((25*(Z12%26+-1==M==0))+1))+((M+15)*(Z12%26+-1==M==0));
                                                        for(N=1;N<10;N++){
                                                            Z14=(Z13/26*((25*(Z13%26+-1==N==0))+1))+((N+5)*(Z13%26+-1==N==0));
                                                            if (Z14==0) {
                                                                printf("%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld\n", A, B, C, D, E, F, G, H, I, J, K, L, M, N);
                                                                exit(0);
                                                            }
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

int main(int argc, char **argv) {
    if (argc == 1) {
        exercise1();
    } else {
        exercise2();
    }
    exit(0);
}
