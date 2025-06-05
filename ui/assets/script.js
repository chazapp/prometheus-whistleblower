function stringToPrometheusMetric(input, value) {
  if (typeof input !== "string") {
    throw new TypeError("Input must be a string.");
  }

  const result = {
    metric: "",
    labels: [],
    value: value
  };

  const metricMatch = input.match(/^([a-zA-Z_:][a-zA-Z0-9_:]*)(\{.*\})?$/);
  if (!metricMatch) {
    throw new Error("Invalid Prometheus metric format.");
  }

  result.metric = metricMatch[1].trim();

  const labelString = metricMatch[2];
  if (labelString) {
    // Remove braces and trim
    const labelsContent = labelString.slice(1, -1).trim();

    // Handle empty label block (e.g., `metric{}`)
    if (labelsContent === "") return result;

    // Match key="value" or key='value'
    const labelPairs = labelsContent.match(/(\w+)=("(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*')/g);
    if (!labelPairs) {
      throw new Error("Invalid label format.");
    }

    result.labels = labelPairs.map(pair => {
      const [label, rawValue] = pair.split('=');
      if (!label || !rawValue) {
        throw new Error(`Malformed label: ${pair}`);
      }

      const value = rawValue.slice(1, -1); // Remove quotes
      return { label, value };
    });
  }

  return result;
}

async function postPrometheusMetric(prometheusMetric) {
  const response = await fetch(
    "/metric", {
      method: "POST",
      body: JSON.stringify(prometheusMetric)
    }
  );
  if (!response.ok) {
    throw new Error(`Error: Response status: ${response.status}`);
  }
}

function in_valueToInt(value) {
  const intValue = parseInt(value, 10);

  if (isNaN(intValue)) {
    throw new Error(`Cannot parse "${value}" to an integer`);
  }

  return intValue;
}


async function submitMetric() {
  const in_metric = document.getElementById("in-metric").value;
  const in_value = document.getElementById("in-value").value;

  try {
    const value = in_valueToInt(in_value);
    const prometheusMetric = stringToPrometheusMetric(in_metric, value);
    console.log(`Sending: ${JSON.stringify(prometheusMetric)}`);
    await postPrometheusMetric(prometheusMetric)
    await loadTable();
  } catch (error) {
    document.getElementById("prom-input-error").textContent = error.message;
  }
}

async function getMetrics() {
  const response = await fetch(
    "/metrics/json"
  );
  if (!response.ok) {
    throw new Error(`Error while fetching metrics: ${response.status}`)
  }
  jsonData = await response.json()
  return jsonData.metrics
}

function jsonLabelsToPromLabelsString(labels) {
  let promLabels = [];

  labels.forEach((kv) => {
    const { label, value } = kv;
    promLabels.push(`${label}="${value}"`)
  })
  return promLabels.join(" ")
}

async function deleteMetric(id) {
  console.log(`Delete item ID:${id}`)
  const response = await fetch(`/metric/${id}`, {
    method: "DELETE"
  })
  if (!response.ok) {
    document.getElementById("metrics-table-error").innerHTML = `Delete error: Response status: ${response.status}`;
    return;
  }
  document.getElementById("metrics-table-error").innerHTML = "";
  await loadTable();
}

async function loadTable() {
  const metrics = await getMetrics();
  console.log(metrics);
  const table = document.getElementById("metrics-table-body");
  table.innerHTML = "";
  await metrics.forEach((metric) => {
    let row = table.insertRow();
    let metricName = row.insertCell(0);
    let labels = row.insertCell(1)
    let value = row.insertCell(2)
    let deleteButton = row.insertCell(3);

    metricName.innerHTML = metric.metric;
    labels.innerHTML = `<code>${jsonLabelsToPromLabelsString(metric.labels)}</code>`
    value.innerHTML = metric.value
    deleteButton.innerHTML = `<button onclick="deleteMetric(${metric.id})">🗑️</button>`
  })
}

loadTable();