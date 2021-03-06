const getData = async () => {
	return fetch("data.json").then((response) => response.json());
};

const MAX_HUE = 120;
const LOWEST = [87, 187, 138]; // '#57BB8A';
const MIDDLE = [255, 214, 101]; // '#FFD665';
const HIGHEST = [230, 123, 114]; // '#E67B72';

const getGradientColor = (color1, color2, min, max, val) => {
	const perc = (val - min) / (max - min);

	const w1 = perc;
	const red = Math.round(color1[0] + (color2[0] - color1[0]) * w1);
	const green = Math.round(color1[1] + (color2[1] - color1[1]) * w1);
	const yellow = Math.round(color1[2] + (color2[2] - color1[2]) * w1);

	return [red, green, yellow];
};

const rgbToCSS = (arr) => `rgb(${arr[0]}, ${arr[1]}, ${arr[2]})`;

const getColor = (val, min, max, percentileMiddle) => {
	if (val === min) {
		return rgbToCSS(LOWEST);
	} else if (val === percentileMiddle) {
		return rgbToCSS(MIDDLE);
	} else if (val === max) {
		return rgbToCSS(HIGHEST);
	} else if (val < percentileMiddle) {
		return rgbToCSS(
			getGradientColor(LOWEST, MIDDLE, min, percentileMiddle, val)
		);
	}

	return rgbToCSS(
		getGradientColor(MIDDLE, HIGHEST, percentileMiddle, max, val)
	);
};

const main = async () => {
	const mainDiv = document.querySelector("#main");

	await getData().then((data) => {
		const max = data.max;
		const min = data.min;
		const percentileMiddle = data.percentileMiddle;
		const values = data.data;

		const elements = [];

		for (let d of values) {
			const date = new Date(d.date).toLocaleDateString();

			elements.push(`<div class="data-row">
				<div class="date-col">${date}</div>
				<div class="data-col">${d.intake
					.map((v) => {
						const c = getColor(v, min, max, percentileMiddle);
						return `<div class="entry" style="--col: ${c}">${v}</div>`;
					})
					.join("")}</div>
			</div>`);
		}

		mainDiv.innerHTML += elements.join("");
	});
};

window.onload = () => main();
