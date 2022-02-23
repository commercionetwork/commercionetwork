# VBR
The `vbr` (*Validator Block Rewards*) module allows validators to get a recurrent reward each 
time a new block is proposed, even if such a block does not contain any transaction.  


## How it works
During genesis, a pool of tokens is created and it's filled with a decided amount of tokens.  
From this point in time, each time a new block is proposed a small amount ok tokens will be sent to the block proposer.

The amount of tokens that will be sent to the validator that proposes the block is computed using the formula
described below. 

### Definitions 
Block produced in a year, considering a block interval of 5 seconds and a year length of 365 days and 8 hours:  
$B = 365.25 \times 24 \times 60 \times \frac{60}{5} = 6,311,520$ 

Initial rewards pool amount:  
$P = 12,500,000$

Voting power of the $n$-th validator at block $b$:  
$VP_{(n, b)} = \frac{S_{(n,b)}}{S_{b}}$

Where
- $S_{(n,b)}$ is the amount of bonded tokens of the $n$-th validator at block $b$   
- $S_{b}$ is the total bonded tokens amount at block $b$ 

Yearly pool reward amount for the $y$-th year:    
$P_y = \begin{cases} P &\text{if } y = 0 \\ P_{y - 1} - \sum_{n, b}{R_{n, b}} &\text{if } y > 0 \end{cases}$

Number of active validators  
$V$

Percentage of active validators towards the target of 100 validators:  
$V_\% = \begin{cases} \frac{V}{100}  &\text{if } V \le 100 \\ 1 &\text{if } V > 100 \end{cases}$

### Formulas
Let's define the maximum reward for the $y$-th year as
  
$$R_y = P_y \times 20\% \times V_\% = P_y \times 20\% \times \frac{V}{100}$$

Define also the maximum limit that each validator can get in a year as 

$$RL_{(y,n)} = \frac{R_y}{100}$$

We must take into consideration that validators with higher voting power must have a lower reward per block
as they will validate more blocks.  
On the other hand, validators with lower voting power will have a higher reward 
per block as they will validate fewer blocks. This results in the given formula

$$R_{(n,b)} = \begin{cases} \frac{RL_{(y,n)}}{B} \times \frac{1}{VP_{(n,b)}} &\text{if } VP_{(n,b)} > \frac{1}{V} \\ \frac{RL_{(y,n)}}{B} \times V &\text{if } VP_{(n,b)} \le \frac{1}{V} \end{cases}$$

