{
  'targets': [
    {
      'target_name': 'nodego',
      'sources': [
        'lib/libgo.h', # auto-generated
        'src/addon.cpp',
        'src/callbacks.cpp',
        'src/promise.cpp',
        'src/value.cpp'
      ],
      'libraries': [
        '../lib/libgo.a' # auto-generated
      ]
    }
  ]
}
