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

    let (part_2, part_1) = readings
        .iter()
        .map(|reading| predict_previous_and_next_reading(reading))
        .reduce(|acc, (previous, next)| (acc.0 + previous, acc.1 + next))
        .unwrap();

    println!("Part 1: {}", part_1);
    println!("Part 2: {}", part_2);
}

fn predict_previous_and_next_reading(previous_readings: &Vec<isize>) -> (isize, isize) {
    if previous_readings.iter().all(|reading| *reading == 0) {
        return (0, 0);
    }

    let (next_sequence_previous_reading, next_sequence_next_reading) =
        predict_previous_and_next_reading(
            &previous_readings
                .windows(2)
                .map(|window| window[1] - window[0])
                .collect(),
        );

    return (
        previous_readings.first().unwrap() - next_sequence_previous_reading,
        next_sequence_next_reading + previous_readings.last().unwrap(),
    );
}

fn read_lines() -> io::Result<io::Lines<io::BufReader<File>>> {
    let filename: String = env::args().skip(1).next().expect("Missing file path");
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}
