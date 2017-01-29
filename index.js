(function () {
  'use strict';

  const WIDTH = 900;
  const HEIGHT = 400;
  const API_URL = 'http://192.168.0.167:8080'
  const paragraph = document.querySelector('span[data-action="voice-recorded"]');
  const signal = document.querySelector('img[data-action="recording-signal"]');
  const canvas = document.getElementById('playground');
  const ctx = canvas.getContext('2d');
  let imageSize = 100;
  let imagesPos = {};

  function init (width, height) {
    // configure canvas setup
    canvas.height = height;
    canvas.width = width;
    canvas.style.background = 'gray';

    // setup areas we should be place images
    createImagePosition(width, height, 3);
    for (var key in imagesPos) {
      drawBorder(imagesPos[key].x, imagesPos[key].y, imageSize, imageSize);
    }

    // setup voice recognition settings
    const BACKGROUND = 'background';
    const MIDDLE = 'middle';
    const LEFT = 'left';
    const RIGHT = 'right';

    const commands = {
      'use a :what as background': put(BACKGROUND),
      'put a :what in the middle': put(MIDDLE),
      'put a :what to the left': put(LEFT),
      'put a :what to the right': put(RIGHT)
    };

    annyang.addCommands(commands);
    annyang.addCallback('error',  function(result) {
      console.log('Ups!!', result)
    })
    annyang.addCallback('result',  recordedString);
    annyang.addCallback('resultNoMatch',  recordedString);
    annyang.addCallback('start',  recordingSignal);

    annyang.start();
  };

  function recordingSignal () {
    signal.style.display = 'block'
  };

  function recordedString (msg) {
    paragraph.innerHTML = msg[0];
  };

  function put2 (where, what) {
    console.log('\\o/', where, what)
    fetch(`${API_URL}/image?what=${what}`)
    .then((res) => {
      return res.text()
    })
    .then((image) => {
      setImagePosition(where, image) 
    })
  };

  function put (where) {
    return what => put2(where, what);
  };

  function drawBorder (x, y, width, height) {
    ctx.lineWidth = 5
    ctx.strokeRect(x, y, width, height)
    ctx.translate(50, 50)
    ctx.setTransform(1, 0, 0, 1, 0, 0)
  };

  function drawImage (url, x, y) {
    var image = new Image()
    image.onload = function () {
      ctx.drawImage(image, x, y, 100, 100)
    }
    image.src = url
  };

  function setImagePosition (where, what) {
    if(where === 'background') {
      canvas.style.background = `url(${what})`
      return;
    }

    drawImage(what, imagesPos[`image${where}`].x, imagesPos[`image${where}`].y)
  };

  function createImagePosition (canvasWidth, canvasHeight, numPlaceholder) {
    var colummSize = canvasWidth / numPlaceholder
    var margin = 50
    var positions = ['left', 'middle', 'right']

    positions.forEach(function (pos, index) {
      // Avoid to start with zero
      index++
      imagesPos[`image${pos}`] = {
        x: index * (colummSize / 2 + margin),
        y: canvasHeight / 2
      }
    })
  };

  init(WIDTH, HEIGHT)
})()
