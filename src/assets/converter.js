import fs from 'fs';
import path from 'path';

// Функция для конвертации шрифта в Base64
function convertFontToBase64(fontPath) {
  const font = fs.readFileSync(fontPath);
  return font.toString('base64');
}

// Указываем путь к шрифту
const fontPath = 'fonts/times.ttf';  // Замените на путь к вашему шрифту

// Получаем Base64 строку шрифта
const base64Font = convertFontToBase64(fontPath);

// Сохраняем Base64 строку в файл
const outputPath = 'font-base64.txt';
fs.writeFileSync(outputPath, base64Font, 'utf8');

console.log(`Base64 шрифта успешно сохранен в файл: ${outputPath}`);