use std::env;
use std::fs::File;
use std::io::{self, BufRead};

fn main() {
    let lines = read_lines().expect("Unable to read lines");
    let readings: Vec<Vec<isize>> = lines
        .map(|line| {
            line.unwrap()
                .split(" ")
                .map(|number| number.parse().unwrap())
                .collect()
        })
        .collect();

    let part_1: isize = readings
        .iter()
        .map(|reading| predict_next_reading(reading))
        .sum();
    println!("Part 1: {}", part_1);
}

fn predict_next_reading(previous_readings: &Vec<isize>) -> isize {
    if previous_readings.iter().all(|reading| *reading == 0) {
        return 0;
    }

    return *previous_readings.last().unwrap()
        + predict_next_reading(
            &previous_readings
                .windows(2)
                .map(|window| window[1] - window[0])
                .collect(),
        );
}

fn read_lines() -> io::Result<io::Lines<io::BufReader<File>>> {
    let filename: String = env::args().skip(1).next().expect("Missing file path");
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}
