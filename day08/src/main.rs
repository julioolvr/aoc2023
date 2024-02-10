use anyhow::anyhow;
use num_integer::Integer;
use petgraph::graph::DiGraph;
use petgraph::visit::{EdgeRef, IntoNodeIdentifiers};
use petgraph::Direction::Outgoing;
use regex::Regex;
use std::env;
use std::fs::File;
use std::io::{self, BufRead};

fn main() -> anyhow::Result<()> {
    let mut lines = read_lines().expect("Unable to read lines");
    let directions = lines
        .next()
        .expect("Could not read directions line")?
        .chars()
        .map(Direction::try_from)
        .map(|result| result.unwrap())
        .collect::<Vec<Direction>>();

    // Skip blank line
    lines.next();

    let mut graph = DiGraph::<String, Direction>::new();
    let node_regex = Regex::new(r"([A-Z0-9]{3}) = \(([A-Z0-9]{3}), ([A-Z0-9]{3})\)")?;

    for line in lines {
        let line = line.unwrap();
        let captures = node_regex.captures(&line).expect("Invalid line");
        let source = captures.get(1).unwrap().as_str().to_string();
        let left = captures.get(2).unwrap().as_str().to_string();
        let right = captures.get(3).unwrap().as_str().to_string();

        let source = graph
            .node_identifiers()
            .find(|id| graph[*id] == source)
            .unwrap_or_else(|| graph.add_node(source));
        let left = graph
            .node_identifiers()
            .find(|id| graph[*id] == left)
            .unwrap_or_else(|| graph.add_node(left));
        let right = graph
            .node_identifiers()
            .find(|id| graph[*id] == right)
            .unwrap_or_else(|| graph.add_node(right));

        graph.add_edge(source, left, Direction::Left);
        graph.add_edge(source, right, Direction::Right);
    }

    let mut current = graph
        .node_indices()
        .find(|index| graph[*index] == "AAA")
        .unwrap();
    let zzz = graph
        .node_indices()
        .find(|index| graph[*index] == "ZZZ")
        .unwrap();
    let mut steps = 0;

    for direction in directions.iter().cycle() {
        steps += 1;
        current = graph
            .edges_directed(current, Outgoing)
            .find(|edge| edge.weight() == direction)
            .unwrap()
            .target();
        if current == zzz {
            break;
        }
    }

    println!("Part 1: {}", steps);

    let starting_points: Vec<_> = graph
        .node_indices()
        .filter(|index| graph[*index].ends_with("A"))
        .collect();

    let cycle_lengths: Vec<usize> = starting_points
        .iter()
        .map(|starting_point| {
            let mut current = *starting_point;
            let mut cycle_length = 0;

            for direction in directions.iter().cycle() {
                cycle_length += 1;

                current = graph
                    .edges_directed(current, Outgoing)
                    .find(|edge| edge.weight() == direction)
                    .unwrap()
                    .target();

                if graph[current].ends_with("Z") {
                    break;
                }
            }

            return cycle_length;
        })
        .collect();

    let steps = cycle_lengths.iter().fold(1, |a, b| a.lcm(b));

    println!("Part 2: {}", steps);
    Ok(())
}

#[derive(Debug, Clone, PartialEq)]
enum Direction {
    Left,
    Right,
}

impl TryFrom<char> for Direction {
    type Error = anyhow::Error;

    fn try_from(value: char) -> Result<Self, Self::Error> {
        match value {
            'L' => Ok(Direction::Left),
            'R' => Ok(Direction::Right),
            _ => Err(anyhow!("Invalid direction")),
        }
    }
}

fn read_lines() -> io::Result<io::Lines<io::BufReader<File>>> {
    let filename: String = env::args().skip(1).next().expect("Missing file path");
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}
