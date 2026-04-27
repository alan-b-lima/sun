# Sun

Sun is a small game based on the ability to cast spells that are multi-stack pushdown automata, which are equivalent to Turing Machines.

## Spells

Spells are multi-stack pushdown automata, technically the same class as Turing Machines, but easier to work with. A spell has an energy level that will be used throughout its casting, if the energy fully depletes, the spell dies. A spell may interacts with the world, absorbing and releasing atoms, and can also move around.

The spells are written in a their own language, the [Sun Spell Specification](./internal/spell/README.md) defines this language.


