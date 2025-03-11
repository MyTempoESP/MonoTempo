# (fork) assertutils.h
A tiny C logging library which also implements exceptions and assertions

This fork is for internal use within MyTempo

![Screenshot_2024-06-22_03-26-05](https://github.com/TsukiGva2/assertutils/assets/71048858/b453759f-352d-48eb-a297-d4a779a07fd3)

# Objective

This project started out at work when i needed a quick and simple logging tool, without the need for excessive configuration and
setup. After some thinkering i came up with this implementation, which made me really happy.

I also must admit that i'm not really proficient at C programming AT ALL, this project will be kept really simple and minimal,
with only some portability tweaks here and there (PRs welcome! 😍)

# Building instructions

The code is really, really simple, you just need to run one of the following:

    $ make

Running `make` will simply build the library as a .so file inside the `lib/` directory, it will also generate the proper folder
structure and headers at `include/`

    # make install

`make install` will automatically compile and try to copy the header files to `/usr/local/include/assertutils/` and the library to
`/usr/local/lib/assertutils/`. You can specify an alternate prefix by passing the DISTDIR option like:

    $ make install DISTDIR=some_directory

The directory will be created automatically, with the correct folder structure (mkdir -p magic).

    $ make example

Running `make example` will compile assertutils along with the examples for you to play around, they are stored in the `bin/` folder,
currently, there are only 3 basic examples: fail.c, success.c and art.c, feel free to play around and explore some interesting use
cases, do consider opening a pull request (or an issue) if you find weird/unexpected behaviour.

# Platform support

This project is a work in progress, the current version is only a simple prototype, I do intend to correctly handle different platforms at a later stage,
but that requires some things to be done first, which takes us to the next section:

## TODOs (PRs welcome)

- [X] Improve the build system (don't look at the Makefile, it's frightening)
- [ ] Reduce the macro abuse in assertutils.h
- [ ] Plan a better error checking model
- [ ] Add tests
- [ ] Make the code more portable (the only things in the way right now really are the colors)
