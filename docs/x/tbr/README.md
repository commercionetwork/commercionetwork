# TBR
The `tbr` (*Transaction Block Rewards*) module allows validators to get a recurrent reward each 
time a new block is proposed, even if such a block does not contain any transaction.  


## How it works
During genesis, a pool of tokens is created and it's filled with a decided amount of tokens.  
From this point in time, each time a new block is proposed a small amount ok tokens will be sent to the block proposer.

The amount of tokens that will be sent to the validator that proposes the block is computed using the following formula:

$$ R_{n,V} = \frac{25000}{365 \times 60 \times 24 \times 12} \times \frac{100}{V} \times \frac{Stake}{TotalStake} $$

Where: 
- $R_{n,V}$ is the reward that will be given to the $n$ th validator considering a set of a total of $V$ validators
- $V$ is the total number of validators present
- $Stake$ is the amount of tokens that has been staked from the $n$ th validator
- $TotalStake$ is the total amount of stake counting all the validators currently active

As an example, considering a set of 100 validators all having the same amount of voting power this would result 
in each one earning up to `25.000` tokens/year. 

However, as you can easily understand from the $\frac{100}{V}$ part of the above formula, if we will have a
validators set which is greater than 100 elements, the amount of tokens that a validator can earn will be lower.  

This is due on purpose in order to make sure that the maximum quantity of tokens that is distributed each year using 
this system will never exceed the 2.5 million tokens cap.