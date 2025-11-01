# Player actions


- Register {username - hashedPassword} uses player resource
- Login    {username - hashedPassword} uses player resource
- Logout   {username - hashedPassword} uses player resource

- BuyBooster       {username - hashedPassword} uses player and card resources
- SaveDeck         {username - hashedPassword - Deck} uses player resource
- CreateTrade  {username - hashedPassword - giveAwayCard} uses Player and Trade Resource
- GetTradableCards {username - hadhedPassword} - uses Player and Trade Resource
- AcceptOffer {username - hashedPassword - TradeId} uses Player and Trade Resources
- SugestTrade {username - hashedPassword - TradeId} uses Player and Trade Resources
- RejectTrade {username - hashedPassword - TradeId} uses Player and Trade Resources

-> A player creates a Trade offer, then another player sugest a card to trade, while this cards are in tradezone 
remove both from player's decks


- Enqueue
- PlayCard
- SkipTurn
- Surrender
