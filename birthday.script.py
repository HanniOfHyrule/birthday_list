import csv
from datetime import datetime
import vobject
import json

def extract_birthdays(csv_filepath):
    birthdays = []
    with open(csv_filepath, mode='r', encoding='utf-8') as file:
        reader = csv.DictReader(file)
        for row in reader:
            name = f"{row['First Name']} {row['Last Name']}".strip()
            birthday = row['Birthday']
            if name and birthday:
                try:
                    datetime.strptime(birthday, '%Y-%m-%d')
                    birthdays.append((name, birthday))
                except ValueError:
                    continue
    return birthdays

def generate_caldav(birthdays, output_filepath):
    calendar = vobject.iCalendar()
    for name, birthday in birthdays:
        event = calendar.add('vevent')
        event.add('summary').value = f"Birthday: {name}"
        
        # Parse the birthday date
        birthday_date = datetime.strptime(birthday, '%Y-%m-%d').date()
        
        # Set start date (all-day event)
        event.add('dtstart').value = birthday_date
        
        # Add recurrence rule for yearly repetition
        event.add('rrule').value = f"FREQ=YEARLY;UNTIL=20801231T235959Z"
        
        # Add description
        event.add('description').value = f"Celebrate {name}'s birthday!"
    
    with open(output_filepath, mode='w', encoding='utf-8') as file:
        file.write(calendar.serialize())

if __name__ == "__main__":
    csv_filepath = "/Users/hannilieber/Development-privat/Birthday_script/contacts.csv"
    output_filepath = "/Users/hannilieber/Development-privat/Birthday_script/birthdays.ics"
    
    birthdays = extract_birthdays(csv_filepath)
    generate_caldav(birthdays, output_filepath)
    print(f"CalDAV file generated at: {output_filepath}")

def save_birthdays_to_json(birthdays, output_filepath):
    # Convert the list of tuples into a list of dictionaries
    birthdays_data = [{"name": name, "birthday": birthday} for name, birthday in birthdays]
    
    # Write the data to a JSON file
    with open(output_filepath, mode='w', encoding='utf-8') as file:
        json.dump(birthdays_data, file, indent=4, ensure_ascii=False)

if __name__ == "__main__":
    csv_filepath = "/Users/hannilieber/Development-privat/Birthday_script/contacts.csv"
    json_output_filepath = "/Users/hannilieber/Development-privat/Birthday_script/birthdays.json"
    
    birthdays = extract_birthdays(csv_filepath)
    save_birthdays_to_json(birthdays, json_output_filepath)
    print(f"JSON file generated at: {json_output_filepath}")