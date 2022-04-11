<!--
order: 0
title: Vbr Overview
parent:
  title: "vbr"
-->

# `Vbr`

## Abstract

The `vbr` (*Validator Block Rewards*) module allows validators to get a recurrent reward each 
time a new block is proposed, even if such a block does not contain any transaction.  


## How it works
During genesis, a pool of tokens is created and it's filled with a decided amount of tokens.  
From this point in `epoch` duration timing, whenever an epoch ends an amount of tokens will be sent to the active validators.

The amount of tokens that will be sent to the validators that proposed the blocks during the epoch is computed using the formula
described below. 

### Definitions 
Block produced in a year, considering a block interval of 5 seconds and a year length of 365 days and 8 hours:  
$B = 365.25 \times 24 \times 60 \times \frac{60}{5} = 6,311,520$ 

Initial rewards pool amount:   
$P = 12,500,000$

Number of active validators  
$V$

Annual distribution for the $n$-th validator:  
$AD_{n} = S_{n} * 0.5 * \frac{V}{100}$

Where
- $S_{n}$ is the amount of bonded tokens of the $n$-th validator   

Possible Epoch duration:  
$E = \begin{cases} 60s &\text minute \\ minute*60 &\text hour \\ hour*24 &\text day \\ hour*24*7 &\text week \\ hour*24*30 &\text month \end{cases}$

Yearly pool reward amount for the $y$-th year:    
$P_y = \begin{cases} P &\text{if } y = 0 \\ P_{y - 1} - \sum_{n, e}{R_{n, e}} &\text{if } y > 0 \end{cases}$

Where
- $R_{n, e}$ is the reward of the $n$-th validator at the end of the $e$-th epoch. 


### Formula
Let's define the maximum reward for the $e$-th epoch as
  
$$R_e = \sum_{v}\frac{AD_{n}}{E}$$


## Contents

1. **[State](01_state.md)**
2. **[Messages](02_messages.md)**
   - [Increment Block Rewards Pool](02_messages.md#Increment-block-rewards-pool)
   - [Set Parameters](02_messages.md#Set-parameters)
3. **[Events](03_events.md)**
   - [Handlers](03_events.md#handlers)
4. **[Parameters](04_params.md)**
5. **[Client](05_client.md)**