package app.snapshot_bitcake;

import java.io.Serializable;

public class Pair implements Serializable {
    public int first;
    public int second;

    public Pair(int first, int second) {
        this.first = first;
        this.second = second;
    }

    @Override
    public String toString() {
        return "(" + first + ", " + second + ")";
    }
}
