{
  'targets': [
    {
      'target_name': 'nodego',
      'sources': [
        'lib/libgo.h', # auto-generated
        'src/addon.cpp'
      ],
      'libraries': [
        '../lib/libgo.a' # auto-generated
      ]
    }
  ]
}
