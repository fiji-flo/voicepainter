/* global annyang */
'use strict';

const BACKGROUND = 'background';
const MIDDLE = 'middle';
const LEFT = 'left';
const RIGHT = 'right';

const commands = {
  'use a :what as background': put(BACKGROUND),
  'put a :what in the middle': put(MIDDLE),
  'put a :what to the left': put(LEFT),
  'put a :what to the right': put(RIGHT),
};

function put2(where, what) {
  console.log('\\o/')
  const p = document.createElement('p');
  p.textContent = `${what} => ${where}`;
  document.body.appendChild(p);
}

function put(where) {
  return what => put2(where, what);
}

annyang.addCommands(commands);
annyang.start();
