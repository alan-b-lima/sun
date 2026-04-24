# Sun

Sun is a automaton game.

## Atoms

The world of the game is made of atoms, they can be water, fire, air, sand, salt, dirt, rock, ect. Also, there are functional atoms, they are called spells, they can detect their surrounding and do something.

## Spells

Spells are multi-stack pushdown automata, technically the same class as Turing Machines, but easier to work with in my opinion. A spell has an energy level that will be used throughout its casting, if the energy fully deplets, the spell dies. A spell may attract, repel, or combine atoms of the world, a spell may also generate new spells.
