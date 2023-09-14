package io.github.netrixframework;

import java.util.concurrent.locks.ReentrantLock;

class Counter {
    private int value = 0;

    private final ReentrantLock lock = new ReentrantLock();

    public Counter() {}

    public int getValue() {
        int result = 0;
        lock.lock();
        result = value;
        lock.unlock();
        return result;
    }

    public void incr() {
        lock.lock();
        value = value + 1;
        lock.unlock();
    }
}
