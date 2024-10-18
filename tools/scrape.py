#!/usr/bin/env python3

import requests

def proj_zero_three():
  resp = requests.get("https://projectzerothree.info/api.php?format=json")
  return resp

def get_stores(dist: int):
  """
  Fetches all stores given an area and distance, -1 if all stores are needed.

  References:
  https://www.7eleven.com.au/store-locator.html -- base URL
  https://www.7eleven.com.au/storelocator-retail/mulesoft/stores?dist=10                                  # need location or none at all
  https://www.7eleven.com.au/storelocator-retail/mulesoft/stores?lat=-33.8688197&long=151.2092955&dist=10 # need to specify location or none at all
  """
  if dist == -1: resp = requests.get("https://www.7eleven.com.au/storelocator-retail/mulesoft/stores")
  else: resp = requests.get(f"https://www.7eleven.com.au/storelocator-retail/mulesoft/stores?dist={dist}")
  return resp

def get_fuel_price(store_no: int) -> requests.models.Response:
  """
  Targets a specific store and fetches the fuel price.
  Needs to also be a fuel store.

  References:
  https://www.7eleven.com.au/storelocator-retail/mulesoft/fuelPrices?storeNo=2362
  """
  target_url = f"https://www.7eleven.com.au/storelocator-retail/mulesoft/fuelPrices?storeNo={store_no}"
  target_store = requests.get(target_url)
  return target_store

def translate_ean(ean: int) -> str:
  translate = {
    52: 'Special Unleaded 91',
    53: 'Special Diesel',
    57: 'Special E10',
    56: 'Supreme+ 98',
    55: 'Extra 95',
  }
  return translate[ean]

def get_all_stores():
  return get_stores(10000).json()['stores']

def get_all_fuel_prices():
  all_stores = get_all_stores()

print(len(get_all_stores()))
