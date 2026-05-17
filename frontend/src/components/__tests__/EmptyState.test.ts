import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import EmptyState from '../EmptyState.vue'

describe('EmptyState', () => {
  it('should render default description for no-data type', () => {
    const wrapper = mount(EmptyState, {
      props: { type: 'no-data' },
    })
    expect(wrapper.text()).toContain('暂无数据，快来创建第一个吧')
  })

  it('should render default description for no-result type', () => {
    const wrapper = mount(EmptyState, {
      props: { type: 'no-result' },
    })
    expect(wrapper.text()).toContain('未找到匹配的结果')
  })

  it('should render custom description', () => {
    const wrapper = mount(EmptyState, {
      props: { description: '暂无记录' },
    })
    expect(wrapper.text()).toContain('暂无记录')
  })

  it('no-result type should show clear filter link', () => {
    const wrapper = mount(EmptyState, {
      props: { type: 'no-result' },
    })
    expect(wrapper.text()).toContain('清除筛选条件')
  })

  it('no-data type should not show clear filter link', () => {
    const wrapper = mount(EmptyState, {
      props: { type: 'no-data' },
    })
    expect(wrapper.find('.clear-link').exists()).toBe(false)
  })

  it('clicking clear link should emit clear event', async () => {
    const wrapper = mount(EmptyState, {
      props: { type: 'no-result' },
    })
    await wrapper.find('.clear-link').trigger('click')
    expect(wrapper.emitted('clear')).toBeTruthy()
  })

  it('no-data type should render action slot content', () => {
    const wrapper = mount(EmptyState, {
      props: { type: 'no-data' },
      slots: {
        action: '<button>创建</button>',
      },
    })
    expect(wrapper.text()).toContain('创建')
  })
})
