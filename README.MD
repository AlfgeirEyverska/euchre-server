# Euchre Server

## ♠️ Euchre Showdown Server API ♣️

**Bring your bot. Battle the odds. Win eternal Euchre glory.**

Welcome to the **Euchre Showdown**, where the cards are dealt and legends are made.

This server powers online Euchre matches between humans and bots. Want to create your own Euchre-playing bot? You’re in the right place. Connect via sockets, speak JSON, and let your code outplay the rest.

---

## 🚀 Quickstart

1. Connect to the server using a socket connection.
2. **Initiate the session by sending a **`Hello`** message**.
3. Listen for incoming JSON `Envelope` messages.
4. Respond using properly structured JSON `ResponseEnvelope` objects.
5. Handle periodic `connectionCheck` messages from the server.
6. Compete and win.

Your bot can be written in **any language** as long as it can:

- Parse incoming JSON messages
- Format and send JSON responses that match the API specs

You can use the python template in the clients directory or get a head start on data encoding and decoding by using the api code directly with a golang bot.

---

## 📦 Envelope Structure

All messages between the server and client use a simple envelope format:

### 🕽️ From Server to Client

```json
{
  "type": "TrickWinnerUpdate",
  "data": { /* actual message data */ },
  "message": "Player 2 won the trick!"
}
```

### 🖙️ From Client to Server

```json
{
  "type": "Play",
  "data": { "response": /* your integer response */ }
}
```

Each message `type` defines the structure of `data`. Here's what you'll see.

---

## 🧠 Message Types

### 👋 `hello`

When you connect, your bot **must** send this initial message:

```json
{
  "type": "hello",
  "data": "hello"
}
```

This lets the server know you're ready to play.

---

### 🔁 `connectionCheck`

The server will periodically check that you're still connected. Respond with:

```json
{
  "type": "connectionCheck",
  "data": "Pong"
}
```

Failing to respond may lead to disconnection.

---

### 🤔 `RequestForResponse`

This is the main decision prompt your bot will receive.
Here is an example sequence of the requests for response that your bot may be prompted with to start a game.

```json
{
    "type": "pickUpOrPass",
    "data": {
        "playerInfo": {
            "playerID": 1,
            "trump": "Not Chosen",
            "flip": "A♠",
            "hand": [
                "K♣",
                "Q♠",
                "K♠",
                "J♠",
                "9♦"
            ]
        },
        "validResponses": {
            "1": "Pass",
            "2": "Pick It Up"
        }
    },
    "message": "Tell the dealer to pick it up or pass."
}
```

```json
{
    "type": "orderOrPass",
    "data": {
        "playerInfo": {
            "playerID": 1,
            "trump": "Not Chosen",
            "flip": "A♠",
            "hand": [
                "K♣",
                "Q♠",
                "K♠",
                "J♠",
                "9♦"
            ]
        },
        "validResponses": {
            "1": "Pass",
            "2": "♥",
            "3": "♦",
            "4": "♣"
        }
    },
    "message": "A♠ Was burried. Order a suit or pass."
}
```

```json
{
    "type": "goItAlone",
    "data": {
        "playerInfo": {
            "playerID": 1,
            "trump": "♥",
            "flip": "A♠",
            "hand": [
                "K♣",
                "Q♠",
                "K♠",
                "J♠",
                "9♦"
            ]
        },
        "validResponses": {
            "1": "Yes",
            "2": "No"
        }
    },
    "message": "Would you like to go it alone?"
}
```

```json
{
    "type": "playCard",
    "data": {
        "playerInfo": {
            "playerID": 1,
            "trump": "♥",
            "flip": "A♠",
            "hand": [
                "K♣",
                "Q♠",
                "K♠",
                "J♠",
                "9♦"
            ]
        },
        "validResponses": {
            "1": "K♣",
            "2": "Q♠",
            "3": "K♠",
            "4": "J♠",
            "5": "9♦"
        }
    },
    "message": "It is your turn. What would you like to play?"
}
```

---

## Updates

---

### 🎴 `DealerUpdate`

```json
{
  "dealer": 2
}
```

Indicates who the dealer is for the round.

---

### 🃏 `PlayJSON`

Your bot plays a card.

```json
{
  "playerID": 1,
  "played": "9♥"
}
```

📝 Card format uses Unicode suits: ♥ (hearts), ♦ (diamonds), ♣ (clubs), ♠ (spades)

---

### 🥇 `WinnerUpdate`

```json
{
  "winner": "Odd Team"
}
```

The final game winner.

---

### 🏆 `TrickWinnerUpdate`

```json
{
  "playerID": 3,
  "action": "won trick"
}
```

---

### 📊 `TrickScoreUpdate`

```json
{
  "evenTrickScore": 2,
  "oddTrickScore": 3
}
```

---

### 🧲 `ScoreUpdate`

```json
{
  "evenScore": 8,
  "oddScore": 7
}
```

---

### 🗺️ `SuitOrdered`

Trump has been selected!

```json
{
  "playerID": 1,
  "action": "Ordered",
  "trump": "♠",
  "goingAlone": false
}
```

---

Your bot should choose from one of the valid responses. **Note:** Response indexing is **1-based**, not 0-based, to better support keyboard users. For example:

- `1` → "pass"
- `1` → "yes", `2` → "no"

Message types like `pickUpOrPass`, `orderOrPass`, `dealerDiscard`, and `goItAlone` will provide prompts like:

- "Tell the dealer to pick it up or pass."
- "You must discard."
- "♣️ Was buried. Order a suit or pass."
- "Would you like to go it alone?" Respond accordingly based on the `validResponses` map included.

---

## 🧰 Example Client Response

Send back a `ResponseEnvelope` like this:

```json
{
  "type": "orderOrPass", // Same string message type received
  "data": {"response": 1} // Integer id of the valid response chosen
}
```

Make sure you send each JSON message as a single newline-delimited string (i.e., terminate with `\n`).

---

## 🤖 Build a Bot

Bring your strategy. Here’s how to make a great bot:

- Parse the JSON envelopes
- Track the hand, trump, flip card, and dealer
- Choose from `validResponses` intelligently
- Send your response within the time limit (if enforced)
- Play your best card, and keep learning from each trick

---
## Test The Bot

To spin up 3 random bots to test your bot against:
`go run clients/main/main.go`
Otherwise, test 4 of your bots against themselves or your bot against itself 4 times.

---

## 🎉 Final Words

I built this server so you can test your Euchre skills—human or AI. Whether you're coding a rule-based strategy or training a neural network, your bot is welcome here.

**So grab your deck. Write your code. And let the best bot win.**

---

🃏 *Need help?*\
Have questions or want to join the competition? Reach out, contribute, or challenge a friend. Let's shuffle up and deal! 🥲

