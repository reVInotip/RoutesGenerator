import folium
import json

# Чтение координат из JSON файла
def read_coords_from_json(filename):
    with open(filename, 'r', encoding='utf-8') as file:
        data = json.load(file)

    if 'coordinates' in data:
        return data['coordinates']
    else:
        raise KeyError("Поле 'coordinates' не найдено в JSON структуре")

# Использование
filename = 'coords.json'
coords = read_coords_from_json(filename)

coords_swapped = [[lon, lat] for lat, lon in coords]

m = folium.Map(location=coords_swapped[0], zoom_start=13)
folium.PolyLine(coords_swapped, color="blue").add_to(m)
m.save("route.html")