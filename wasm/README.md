# Salat WASM - Browser API

WebAssembly version of Salat prayer time calculator for browser usage.

## CDN Usage

### Quick Start

```html
<!DOCTYPE html>
<html>
<head>
    <title>Salat Prayer Times</title>
</head>
<body>
    <div id="prayer-times"></div>
    
    <!-- Load from CDN -->
    <script src="https://github.com/herbras/BelGolang/releases/latest/download/salat.js"></script>
    
    <script>
        // Calculate prayer times for Jakarta
        Salat.prayerTimes(-6.2088, 106.8456, 'Kemenag')
            .then(result => {
                document.getElementById('prayer-times').innerHTML = 
                    `<h3>Prayer Times for Jakarta</h3>
                     <ul>
                         <li>Imsak: ${result.prayers.imsak}</li>
                         <li>Subuh: ${result.prayers.subuh}</li>
                         <li>Dzuhur: ${result.prayers.dzuhur}</li>
                         <li>Ashar: ${result.prayers.ashar}</li>
                         <li>Maghrib: ${result.prayers.maghrib}</li>
                         <li>Isya: ${result.prayers.isya}</li>
                     </ul>
                     <p>Current: ${result.current.emoji} ${result.current.prayer}</p>
                     <p>Next: ${result.next.emoji} ${result.next.prayer} at ${result.next.time}</p>`;
            });
    </script>
</body>
</html>
```

## API Reference

### Methods

#### `Salat.prayerTimes(latitude, longitude, method?)`
Calculate prayer times for given coordinates.

**Parameters:**
- `latitude` (number): Latitude coordinate
- `longitude` (number): Longitude coordinate  
- `method` (string, optional): Calculation method (default: 'Kemenag')

**Methods available:**
- `MWL` - Muslim World League
- `ISNA` - Islamic Society of North America
- `Egypt` - Egyptian General Authority of Survey
- `Makkah` - Umm al-Qura University, Makkah
- `Karachi` - University of Islamic Sciences, Karachi
- `Tehran` - Institute of Geophysics, University of Tehran
- `Kemenag` - Kementerian Agama Republik Indonesia (default)
- `JAKIM` - Jabatan Kemajuan Islam Malaysia

**Returns:** Promise with prayer times object

#### `Salat.version()`
Get library version and build info.

#### `Salat.command(cmd)`
Execute CLI-style commands.

**Examples:**
```javascript
// Get help
await Salat.command('help');

// Get version
await Salat.command('version');

// Calculate prayer times
await Salat.command('prayer -6.2088 106.8456 Kemenag');

// List methods
await Salat.command('methods');
```

## Terminal Integration

### Create Interactive Terminal

```html
<div id="terminal" style="height: 300px; overflow-y: auto; background: #000; color: #0f0; padding: 10px;"></div>

<script>
const terminal = Salat.createTerminal(document.getElementById('terminal'));

// Execute commands
terminal.execute('prayer -6.2088 106.8456');
terminal.execute('version');
terminal.execute('help');
</script>
```

### Chat App Integration

```javascript
// For chat applications
async function handleSalatCommand(userMessage) {
    if (userMessage.startsWith('/salat')) {
        const command = userMessage.replace('/salat', 'prayer');
        const result = await Salat.command(command);
        
        if (result.error) {
            return `❌ ${result.error}`;
        }
        
        return `🕌 Prayer Times:
${result.current.emoji} Current: ${result.current.prayer}
${result.next.emoji} Next: ${result.next.prayer} at ${result.next.time}

📍 Location: ${result.location.latitude}, ${result.location.longitude}
📅 Method: ${result.method}

⏰ Today's Schedule:
🌙 Imsak: ${result.prayers.imsak}
🌅 Subuh: ${result.prayers.subuh}  
☀️ Dzuhur: ${result.prayers.dzuhur}
🌤️ Ashar: ${result.prayers.ashar}
🌇 Maghrib: ${result.prayers.maghrib}
✨ Isya: ${result.prayers.isya}`;
    }
}

// Usage in chat
handleSalatCommand('/salat -6.2088 106.8456 Kemenag')
    .then(response => console.log(response));
```

## Advanced Usage

### Custom Loading

```javascript
// Load with custom options
await Salat.load({
    debug: true,
    wasmURL: 'https://your-cdn.com/salat.wasm'
});
```

### Quick Helpers

```javascript
// Get just prayer times
const times = await Salat.getPrayerTimesForLocation(-6.2088, 106.8456);

// Get current prayer info
const current = await Salat.getCurrentPrayer(-6.2088, 106.8456);
console.log(`Current prayer: ${current.current}, Next: ${current.next}`);
```

## File Structure

When hosting yourself:
```
your-cdn/
├── salat.js        # Main loader script
├── salat.wasm      # WebAssembly module  
└── wasm_exec.js    # Go WASM runtime (auto-loaded)
```

## Browser Support

- Chrome 57+
- Firefox 52+ 
- Safari 11+
- Edge 16+

WebAssembly support required. 