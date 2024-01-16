use anyhow::anyhow;
use std::env;
use std::fs::File;
use std::io::{self, BufRead};
use std::{collections::HashMap, str::FromStr};

fn main() -> anyhow::Result<()> {
    let mut hands = read_lines()
        .expect("Unable to read lines")
        .map(|line| {
            line.map_err(anyhow::Error::from)
                .and_then(|line| line.parse())
        })
        .collect::<Result<Vec<Hand>, _>>()?;

    hands.sort();

    let part_1: usize = hands
        .iter()
        .enumerate()
        .map(|(i, hand)| (i + 1) * hand.bid)
        .sum();
    println!("Part 1: {}", part_1);

    let mut joker_hands = hands.into_iter().map(JokerHand).collect::<Vec<_>>();
    joker_hands.sort();

    let part_2: usize = joker_hands
        .iter()
        .enumerate()
        .map(|(i, hand)| (i + 1) * hand.0.bid)
        .sum();
    println!("Part 2: {}", part_2);

    Ok(())
}

#[derive(PartialEq, Debug, Eq)]
struct Hand {
    cards: Vec<CardLabel>,
    bid: usize,
}

impl Hand {
    fn get_hand_type(&self) -> HandType {
        let counts = self.cards.iter().fold(HashMap::new(), |mut acc, card| {
            *acc.entry(card).or_insert(0) += 1;
            acc
        });

        if counts.len() == 1 {
            return HandType::FiveOfAKind;
        }

        if counts.values().any(|&count| count == 4) {
            return HandType::FourOfAKind;
        }

        if counts.values().any(|&count| count == 3) {
            if counts.values().any(|&count| count == 2) {
                return HandType::FullHouse;
            } else {
                return HandType::ThreeOfAKind;
            }
        }

        if counts.values().filter(|&count| *count == 2).count() == 2 {
            return HandType::TwoPair;
        }

        if counts.values().any(|&count| count == 2) {
            return HandType::OnePair;
        }

        HandType::HighCard
    }

    fn get_joker_hand_type(&self) -> HandType {
        let non_joker_hand_type = self.get_hand_type();
        let joker_hand_type: HandType;

        let mut counts = self.cards.iter().fold(HashMap::new(), |mut acc, card| {
            *acc.entry(card).or_insert(0) += 1;
            acc
        });

        let joker_counts = counts.remove(&CardLabel::Jack).unwrap_or(0);

        if joker_counts == 0 {
            return non_joker_hand_type;
        }

        if joker_counts == 5 {
            return HandType::FiveOfAKind;
        }

        if counts.values().max().unwrap_or(&0) + joker_counts == 5 {
            // If the max non-Joker count + the Joker count is 5, then we have a Five of a Kind
            joker_hand_type = HandType::FiveOfAKind;
        } else if counts.values().max().unwrap_or(&0) + joker_counts == 4 {
            // If the max non-Joker count + the Joker count is 4, then we have a Four of a Kind
            joker_hand_type = HandType::FourOfAKind;
        } else if counts.values().max().unwrap_or(&0) + joker_counts == 3 {
            // If there are 3 jokers or more, we'd have either 4 of a kind or 5 of a kind.
            // So here we have either one or two jokers (the 0 jokers scenario is an early return
            // above).
            if joker_counts == 2 {
                // * If there's two jokers, then the remaining cards must be all different because
                //   otherwise we'd have either four or five of a kind. That means we have
                //   ThreeOfAKind (any of those cards, with the two jokers).
                joker_hand_type = HandType::ThreeOfAKind;
            } else {
                // * If there's one joker, then there must be at least one pair of cards that are
                //   equal to get the total of 3 from this conditional branch. If the remaining two
                //   cards are a pair, we get a FullHouse. Otherwise we have ThreeOfAKind.
                if counts.values().all(|&count| count == 2) {
                    joker_hand_type = HandType::FullHouse;
                } else {
                    joker_hand_type = HandType::ThreeOfAKind;
                }
            }
        } else {
            // The only remaining scenario is that the max count + the number of jokers is 2. There
            // is at least one joker, so all the remaining cards are different. This yields a
            // OnePair by combining the joker with any card
            joker_hand_type = HandType::OnePair;
        }

        return if joker_hand_type > non_joker_hand_type {
            joker_hand_type
        } else {
            non_joker_hand_type
        };
    }
}

impl FromStr for Hand {
    type Err = anyhow::Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let [cards, bid] = s.split_whitespace().collect::<Vec<&str>>()[..] else {
            return Err(anyhow!("Invalid input"));
        };

        let bid = bid.parse::<usize>()?;
        let cards = cards
            .chars()
            .map(|card| match card {
                '2' => Ok(CardLabel::Two),
                '3' => Ok(CardLabel::Three),
                '4' => Ok(CardLabel::Four),
                '5' => Ok(CardLabel::Five),
                '6' => Ok(CardLabel::Six),
                '7' => Ok(CardLabel::Seven),
                '8' => Ok(CardLabel::Eight),
                '9' => Ok(CardLabel::Nine),
                'T' => Ok(CardLabel::Ten),
                'J' => Ok(CardLabel::Jack),
                'Q' => Ok(CardLabel::Queen),
                'K' => Ok(CardLabel::King),
                'A' => Ok(CardLabel::Ace),
                _ => Err(anyhow!("Invalid card label")),
            })
            .collect::<Result<Vec<CardLabel>, anyhow::Error>>()?;

        Ok(Hand { bid, cards })
    }
}

impl PartialOrd for Hand {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        use std::cmp::Ordering::*;

        let self_hand_type = self.get_hand_type();
        let other_hand_type = other.get_hand_type();

        if self_hand_type > other_hand_type {
            return Some(Greater);
        } else if self_hand_type < other_hand_type {
            return Some(Less);
        } else {
            return self
                .cards
                .iter()
                .zip(other.cards.iter())
                .map(|(a, b)| a.cmp(b))
                .filter(|result| *result != Equal)
                .next();
        }
    }
}

impl Ord for Hand {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        self.partial_cmp(other).unwrap()
    }
}

#[derive(PartialEq, Eq)]
struct JokerHand(Hand);

impl PartialOrd for JokerHand {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        use std::cmp::Ordering::*;

        let self_hand_type = self.0.get_joker_hand_type();
        let other_hand_type = other.0.get_joker_hand_type();

        if self_hand_type > other_hand_type {
            return Some(Greater);
        } else if self_hand_type < other_hand_type {
            return Some(Less);
        } else {
            return self
                .0
                .cards
                .iter()
                .zip(other.0.cards.iter())
                .map(|(a, b)| match (a, b) {
                    // Manually compare the Jack / Joker case so that it's smaller than anything
                    // else without having to reimplement PartialOrd for CardLabel
                    (CardLabel::Jack, CardLabel::Jack) => Equal,
                    (CardLabel::Jack, _) => Less,
                    (_, CardLabel::Jack) => Greater,
                    _ => a.cmp(b),
                })
                .filter(|result| *result != Equal)
                .next();
        }
    }
}

impl Ord for JokerHand {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        self.partial_cmp(other).unwrap()
    }
}

#[derive(PartialEq, Eq, PartialOrd, Ord, Hash, Debug)]
enum CardLabel {
    Two,
    Three,
    Four,
    Five,
    Six,
    Seven,
    Eight,
    Nine,
    Ten,
    Jack,
    Queen,
    King,
    Ace,
}

#[derive(PartialEq, Debug)]
enum HandType {
    FiveOfAKind,
    FourOfAKind,
    FullHouse,
    ThreeOfAKind,
    TwoPair,
    OnePair,
    HighCard,
}

impl PartialOrd for HandType {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        use std::cmp::Ordering::*;

        match self {
            HandType::FiveOfAKind => match other {
                HandType::FiveOfAKind => Some(Equal),
                _ => Some(Greater),
            },
            HandType::FourOfAKind => match other {
                HandType::FiveOfAKind => Some(Less),
                HandType::FourOfAKind => Some(Equal),
                _ => Some(Greater),
            },
            HandType::FullHouse => match other {
                HandType::FiveOfAKind => Some(Less),
                HandType::FourOfAKind => Some(Less),
                HandType::FullHouse => Some(Equal),
                _ => Some(Greater),
            },
            HandType::ThreeOfAKind => match other {
                HandType::FiveOfAKind => Some(Less),
                HandType::FourOfAKind => Some(Less),
                HandType::FullHouse => Some(Less),
                HandType::ThreeOfAKind => Some(Equal),
                _ => Some(Greater),
            },
            HandType::TwoPair => match other {
                HandType::FiveOfAKind => Some(Less),
                HandType::FourOfAKind => Some(Less),
                HandType::FullHouse => Some(Less),
                HandType::ThreeOfAKind => Some(Less),
                HandType::TwoPair => Some(Equal),
                _ => Some(Greater),
            },
            HandType::OnePair => match other {
                HandType::FiveOfAKind => Some(Less),
                HandType::FourOfAKind => Some(Less),
                HandType::FullHouse => Some(Less),
                HandType::ThreeOfAKind => Some(Less),
                HandType::TwoPair => Some(Less),
                HandType::OnePair => Some(Equal),
                _ => Some(Greater),
            },
            HandType::HighCard => match other {
                HandType::HighCard => Some(Equal),
                _ => Some(Less),
            },
        }
    }
}

fn read_lines() -> io::Result<io::Lines<io::BufReader<File>>> {
    let filename: String = env::args().skip(1).next().expect("Missing file path");
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}
